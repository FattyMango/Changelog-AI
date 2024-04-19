package models


// map[inference_status:map[cost:1.404e-05 runtime_ms:561 tokens_generated:40 tokens_input:12] num_input_tokens:12 num_tokens:40 results:[map[generated_text: Hello! Greetings from Replika! I'm here to help answer your questions, share in conversation, and get to know you better. How has your day been so far?]]]

type BaseResponse struct {
	InferenceStatus map[string]interface{} `json:"inference_status"`
	NumInputTokens  int                    `json:"num_input_tokens"`
	NumTokens       int                    `json:"num_tokens"`
	Results         []Result               `json:"results"`
}

type Result struct {
	GeneratedText string `json:"generated_text"`
}
