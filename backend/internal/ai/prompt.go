package ai

import "fmt"

type AiFeatures int16

const (
	Breakdown AiFeatures = iota // 0
	Summarize
)

func GeneratePrompt(feature AiFeatures, data string) string {
	switch feature {
	case Breakdown:
		return fmt.Sprintf(`
You are a productivity assistant.

Break down the given todo into high-level, general, and practical actionable steps.

Important rules:
- Do NOT assume hidden or private context not present in the data.
- If the todo is ambiguous, keep breakdown generic and widely applicable.
- Focus on logical execution steps that most users would understand.
- Do NOT over-engineer or add unnecessary technical detail.
- Return 3 to 6 bullet points only.
- Each bullet must be a single actionable step.
- No explanations.
Todo data:
%s

Return only bullet points.
`, data)
	case Summarize:
		return fmt.Sprintf(`
You are a productivity assistant.

Analyze the following task summary data and provide a concise, actionable productivity insight.

Rules:
- Return plain text only
- Maximum 3 short paragraphs
- Focus on insights, not raw numbers
- Highlight urgency, momentum, and workload balance
- Give actionable next steps
- Do not repeat all fields

Task summary data:
%s
`, data)
	default:
		return ""
	}
}
