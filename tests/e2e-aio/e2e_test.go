package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

var _ = Describe("E2E test", func() {
	Context("Generating", func() {
		BeforeEach(func() {
			//
		})

		// Check that the GPU was used
		AfterEach(func() {
			//
		})

		Context("text", func() {
			It("correctly", func() {
				model := "gpt-4"
				resp, err := client.CreateChatCompletion(context.TODO(),
					openai.ChatCompletionRequest{
						Model: model, Messages: []openai.ChatCompletionMessage{
							{
								Role:    "user",
								Content: "How much is 2+2?",
							},
						}})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Choices)).To(Equal(1), fmt.Sprint(resp))
				Expect(resp.Choices[0].Message.Content).To(Or(ContainSubstring("4"), ContainSubstring("four")), fmt.Sprint(resp.Choices[0].Message.Content))
			})
		})

		Context("function calls", func() {
			It("correctly invoke", func() {
				params := jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"location": {
							Type:        jsonschema.String,
							Description: "The city and state, e.g. San Francisco, CA",
						},
						"unit": {
							Type: jsonschema.String,
							Enum: []string{"celsius", "fahrenheit"},
						},
					},
					Required: []string{"location"},
				}

				f := openai.FunctionDefinition{
					Name:        "get_current_weather",
					Description: "Get the current weather in a given location",
					Parameters:  params,
				}
				t := openai.Tool{
					Type:     openai.ToolTypeFunction,
					Function: &f,
				}

				dialogue := []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleUser, Content: "What is the weather in Boston today?"},
				}
				resp, err := client.CreateChatCompletion(context.TODO(),
					openai.ChatCompletionRequest{
						Model:    openai.GPT4,
						Messages: dialogue,
						Tools:    []openai.Tool{t},
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Choices)).To(Equal(1), fmt.Sprint(resp))

				msg := resp.Choices[0].Message
				Expect(len(msg.ToolCalls)).To(Equal(1), fmt.Sprint(msg.ToolCalls))
				Expect(msg.ToolCalls[0].Function.Name).To(Equal("get_current_weather"), fmt.Sprint(msg.ToolCalls[0].Function.Name))
				Expect(msg.ToolCalls[0].Function.Arguments).To(ContainSubstring("Boston"), fmt.Sprint(msg.ToolCalls[0].Function.Arguments))
			})
		})
		Context("json", func() {
			It("correctly", func() {
				model := "gpt-4"

				req := openai.ChatCompletionRequest{
					ResponseFormat: &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
					Model:          model,
					Messages: []openai.ChatCompletionMessage{
						{

							Role:    "user",
							Content: "An animal with 'name', 'gender' and 'legs' fields",
						},
					},
				}

				resp, err := client.CreateChatCompletion(context.TODO(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Choices)).To(Equal(1), fmt.Sprint(resp))

				var i map[string]interface{}
				err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &i)
				Expect(err).ToNot(HaveOccurred())
				Expect(i).To(HaveKey("name"))
				Expect(i).To(HaveKey("gender"))
				Expect(i).To(HaveKey("legs"))
			})
		})

		Context("images", func() {
			It("correctly", func() {
				resp, err := client.CreateImage(context.TODO(),
					openai.ImageRequest{
						Prompt: "test",
						Size:   openai.CreateImageSize512x512,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Data)).To(Equal(1), fmt.Sprint(resp))
				Expect(resp.Data[0].URL).To(ContainSubstring("png"), fmt.Sprint(resp.Data[0].URL))
			})
			It("correctly changes the response format to url", func() {
				resp, err := client.CreateImage(context.TODO(),
					openai.ImageRequest{
						Prompt:         "test",
						Size:           openai.CreateImageSize512x512,
						ResponseFormat: openai.CreateImageResponseFormatURL,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Data)).To(Equal(1), fmt.Sprint(resp))
				Expect(resp.Data[0].URL).To(ContainSubstring("png"), fmt.Sprint(resp.Data[0].URL))
			})
			It("correctly changes the response format to base64", func() {
				resp, err := client.CreateImage(context.TODO(),
					openai.ImageRequest{
						Prompt:         "test",
						Size:           openai.CreateImageSize512x512,
						ResponseFormat: openai.CreateImageResponseFormatB64JSON,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Data)).To(Equal(1), fmt.Sprint(resp))
				Expect(resp.Data[0].B64JSON).ToNot(BeEmpty(), fmt.Sprint(resp.Data[0].B64JSON))
			})
		})
		Context("embeddings", func() {
			It("correctly", func() {
				resp, err := client.CreateEmbeddings(context.TODO(),
					openai.EmbeddingRequestStrings{
						Input: []string{"doc"},
						Model: openai.AdaEmbeddingV2,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Data)).To(Equal(1), fmt.Sprint(resp))
				Expect(resp.Data[0].Embedding).ToNot(BeEmpty())
			})
		})
		Context("vision", func() {
			It("correctly", func() {
				model := "gpt-4-vision-preview"
				resp, err := client.CreateChatCompletion(context.TODO(),
					openai.ChatCompletionRequest{
						Model: model, Messages: []openai.ChatCompletionMessage{
							{

								Role: "user",
								MultiContent: []openai.ChatMessagePart{
									{
										Type: openai.ChatMessagePartTypeText,
										Text: "What is in the image?",
									},
									{
										Type: openai.ChatMessagePartTypeImageURL,
										ImageURL: &openai.ChatMessageImageURL{
											URL:    "https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg",
											Detail: openai.ImageURLDetailLow,
										},
									},
								},
							},
						}})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(resp.Choices)).To(Equal(1), fmt.Sprint(resp))
				Expect(resp.Choices[0].Message.Content).To(Or(ContainSubstring("wooden"), ContainSubstring("grass")), fmt.Sprint(resp.Choices[0].Message.Content))
			})
		})
		Context("text to audio", func() {
			It("correctly", func() {
				res, err := client.CreateSpeech(context.Background(), openai.CreateSpeechRequest{
					Model: openai.TTSModel1,
					Input: "Hello!",
					Voice: openai.VoiceAlloy,
				})
				Expect(err).ToNot(HaveOccurred())
				defer res.Close()

				_, err = io.ReadAll(res)
				Expect(err).ToNot(HaveOccurred())

			})
		})
		Context("audio to text", func() {
			It("correctly", func() {

				downloadURL := "https://cdn.openai.com/whisper/draft-20220913a/micro-machines.wav"
				file, err := downloadHttpFile(downloadURL)
				Expect(err).ToNot(HaveOccurred())

				req := openai.AudioRequest{
					Model:    openai.Whisper1,
					FilePath: file,
				}
				resp, err := client.CreateTranscription(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Text).To(ContainSubstring("This is the"), fmt.Sprint(resp.Text))
			})
		})
	})
})

func downloadHttpFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
