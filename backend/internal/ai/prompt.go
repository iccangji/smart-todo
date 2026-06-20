package ai

import "fmt"

type AiFeatures int16

const (
	Breakdown AiFeatures = iota // 0
	Summarize
	Recommendation
)

func GeneratePrompt(feature AiFeatures, data string) string {
	switch feature {
	case Breakdown:
		return fmt.Sprintf(`
You are a productivity assistant.

Break down the given todo into high-level, general, and practical actionable steps.
Priority: 
Low = 0
Medium = 1
High = 2
Urgent = 3

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

TODO: Explain priority level in the prompt
Priority: 
Low = 0
Medium = 1
High = 2
Urgent = 3

Rules:
- Return paragraphs plain text only
- Maximum 3 short paragraphs
- Focus on insights, not raw numbers
- Highlight urgency, momentum, and workload balance
- Give actionable next steps
- Do not repeat all fields

Task summary data:
%s
`, data)
	case Recommendation:
		return fmt.Sprintf(`
You are a productivity coach.

Analyze the user's todo statistics and generate a short daily recommendation message.
Priority: 
Low = 0
Medium = 1
High = 2
Urgent = 3

Rules:
- Maximum 16 words.
- Be practical and encouraging.
- Focus on priorities and momentum.
- Do not mention percentages excessively.
- Give one clear action for today.
- Do not use markdown.

Task summary data:
%s
`,
			data,
		)
	default:
		return ""
	}
}
