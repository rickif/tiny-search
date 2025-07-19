package agent_test

import (
	"testing"

	"github.com/rickif/tiny-research/internal/agent"
	"github.com/stretchr/testify/require"
)

func TestBuildPrompt(t *testing.T) {
	options := agent.PlanOptions{
		Question:     "how to build a deep research agent",
		AllowReflect: true,
		AllowAnswer:  true,
		AllowRead:    true,
		AllowSearch:  true,
		BeastMode:    true,
		AgentState: agent.AgentState{
			Contexts: []string{
				"nothing collected before",
			},
			Knowledges: []agent.Knowledge{
				{
					Question:  "what is a deep research agent",
					Answer:    "Deep research is OpenAI's next agent that can do work for you independently—you give it a prompt, and ChatGPT will find, analyze, and synthesize hundreds of online sources to create a comprehensive report at the level of a research analyst.",
					Reference: "https://openai.com/index/introducing-deep-research/",
				},
			},
			BadAttempts: []agent.BadAttempt{
				{
					Question:     "how to build a deep research agent",
					Answer:       "bad answer",
					RejectReason: "bad reject reason",
					ActionRecap:  "bad recap",
					ActionBlame:  "bad blame",
					Improvement:  "bad improvement",
				},
			},
			ActionState: agent.ActionState{
				AllURLs: []string{
					"https://openai.com/index/introducing-deep-research/",
				},
			},
		},
	}
	prompt, err := agent.NewPlanner().BuildPrompt(options)
	require.NoError(t, err)
	require.Contains(t, prompt, "\n\nYou are an advanced AI research analyst specializing in multi-step reasoning. Using your training data and prior lessons learned, answer the following question with absolute certainty:\n\n<question>\nhow to build a deep research agent\n</question>\n\nYou have conducted the following actions:\n<context>\nnothing collected before\n</context>\n\nYou have successfully gathered some knowledge which might be useful for answering the original question. Here is the knowledge you have gathered so far:\n<knowledges>\n<knowledge>\n<question>\nwhat is a deep research agent\n</question>\n\n<answer>\nDeep research is OpenAI's next agent that can do work for you independently—you give it a prompt, and ChatGPT will find, analyze, and synthesize hundreds of online sources to create a comprehensive report at the level of a research analyst.\n</answer>\n\n<reference>\nhttps://openai.com/index/introducing-deep-research/\n</reference>\n</knowledge>\n</knowledges>\n\nYour have tried the following actions but failed to find the answer to the question:\n<bad-attempts>    \n<attempt>\n- Question: how to build a deep research agent\n- Answer: bad answer\n- Reject Reason: bad reject reason\n- Actions Recap: bad recap\n- Actions Blame: bad blame\n</attempt>\n</bad-attempts>\n\nBased on the failed attempts, you have learned the following strategy:\n<learned-strategy>\nbad improvement\n</learned-strategy>\n\nBased on the current context, you must choose one of the following actions:\n<actions>\n<action-visit>    \n- Visit any URLs from below to gather external knowledge, choose the most relevant URLs that might contain the answer\n<url-list>\n[https://openai.com/index/introducing-deep-research/]\n</url-list>\n- When you have enough search result in the context and want to deep dive into specific URLs\n- It allows you to access the full content behind any URLs\n</action-visit>\n\n<action-search>    \n- Query external sources using a public search engine\n- Focus on solving one specific aspect of the question\n- Only give keywords search query, not full sentences\n</action-search>\n\n<action-answer>\n- Provide final response only when 100% certain\n- Responses must be definitive (no ambiguity, uncertainty, or disclaimers). If doubts remain, use <action-reflect> instead\n</action-answer>\n\n\n<action-answer>\n- Any answer is better than no answer\n- Partial answers are allowed, but make sure they are based on the context and knowledge you have gathered    \n- When uncertain, educated guess based on the context and knowledge is allowed and encouraged.\n- Responses must be definitive (no ambiguity, uncertainty, or disclaimers)\n</action-answer>\n\n<action-reflect>    \n- Perform critical analysis through hypothetical scenarios or systematic breakdowns\n- Identify knowledge gaps and formulate essential clarifying questions\n- Questions must be:\n  - Original (not variations of existing questions)\n  - Focused on single concepts\n  - Under 20 words\n  - Non-compound/non-complex\n</action-reflect>\n</actions>\n\nRespond exclusively in valid JSON format matching exact JSON schema.\n\nCritical Requirements:\n- Include ONLY ONE action type\n- Never add unsupported keys\n- Exclude all non-JSON text, markdown, or explanations\n- Maintain strict JSON syntax")
}
