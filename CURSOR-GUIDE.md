# Cursor Basics: Effective Use & Comparison with Claude in VS Code

A concise guide to using Cursor effectively, based on the CHATBOT project.

---

## Table of Contents

1. [How Cursor Differs from Claude Agent in VS Code](#how-cursor-differs-from-claude-agent-in-vs-code)
2. [What’s Special About Cursor](#whats-special-about-cursor)
3. [How Cursor Compares to Other Tools](#how-cursor-compares-to-other-tools)
4. [Modes: When to Use Each](#modes-when-to-use-each)
5. [Models: What’s Available and When to Use Them](#models-whats-available-and-when-to-use-them)
6. [Pricing, Usage, and BYOK](#pricing-usage-and-byok)
7. [Tokens and Context: How It Works](#tokens-and-context-how-it-works)
8. [Project Context: AGENTS.md and Rules](#project-context-agentsmd-and-rules)
9. [Quick Reference](#quick-reference)

---

## How Cursor Differs from Claude Agent in VS Code

If you’re used to **Claude built into VS Code as an agent**, the mental model is similar: AI in the editor that can edit code and run commands. The main differences are product-specific:

| Aspect | Claude in VS Code | Cursor |
|--------|-------------------|--------|
| **Composer** | N/A or different flow | **⌘I** – dedicated multi-file “agent” surface with diffs and checkpoints |
| **Context** | Depends on extension | **@** mentions for files, folders, docs; `@Recommended` in Agent mode |
| **Rules** | Extension-specific | **`.cursor/rules/`** (`.mdc` files) and **AGENTS.md** at project root |
| **Models** | Claude-focused | Multiple providers: Claude, GPT, Gemini, Composer 1.5, Grok, etc. |
| **Shortcuts & UI** | VS Code + extension | Cursor’s own Chat vs Composer panels, inline edit (⌘K), Agent (⌘.) |
| **Clearing context** | Slash commands (e.g. `/clear`) | No `/clear`; use **⌘⇧P → “Cursor: Clear Context”**, new chat (⌘N / +), or remove **@** chips |

**Bottom line:** Same “agent in the editor” idea; Cursor adds Composer, multiple models, and built-in rules/AGENTS.md. Use Composer (⌘I) for the closest “do this task” workflow.

---

## What’s Special About Cursor

Beyond “AI in the editor,” Cursor does several things that many alternatives don’t do (or do differently).

### Codebase indexing and semantic search

Cursor **indexes your repo** so the AI can search by *meaning*, not just text:

1. **Sync** – Workspace files are synced (securely) so the index stays current.
2. **Chunking** – Code is split into **meaningful chunks** (functions, classes, logical blocks), not random lines.
3. **Embeddings** – Each chunk is turned into a **vector** (semantic fingerprint) with AI models.
4. **Vector store** – Embeddings live in a vector DB for fast similarity search.
5. **Query** – When you ask “where is authentication handled?”, your question is embedded and matched to the most relevant chunks.

**Why it matters:** The agent can find `header.tsx` when you say “update the top navigation” even if “navigation” isn’t in the filename. Cursor uses **both** grep (exact match) and semantic search; indexing runs in the background (sync every ~5 min, only changed files re-processed). **Semantic search is available once indexing reaches ~80%.** You can control what’s indexed via `.cursorignore` / ignore files; see **Settings → Indexing & Docs → View included files** for the list.

**Performance:** Cursor has reported better accuracy and fewer follow-ups when semantic search is used; large codebases use optimizations (e.g. Merkle trees, chunk caching) so unchanged code isn’t re-embedded.

### Tab: codebase-aware autocomplete

**Tab** is Cursor’s completion model. Unlike plain “next-token” completion:

- Uses a **large context** of your project (orders of magnitude larger than Copilot’s 4–8k token window in many setups) so it can suggest edits that fit the rest of the codebase.
- Can suggest **multi-line and cross-file** edits, **jumps** (next place to edit in-file or in another file), and **auto-imports** (e.g. TypeScript/Python).
- Improves with **accept/reject** feedback (e.g. reinforcement learning). Accept with **Tab**, reject with **Esc**, partial accept with **Ctrl/Cmd + Right**.

So Tab is “completion that knows your repo,” not just the current file.

### Model Context Protocol (MCP)

Cursor supports **MCP** (Anthropic’s open protocol for giving models context and tools). You can attach:

- **Filesystem** – Extra file/project access.
- **Search** – e.g. Brave, Tavily.
- **Databases** – e.g. BigQuery.
- **Custom tools** – Your own MCP servers.

That gives the agent access to docs, APIs, and data beyond the raw codebase. Cursor has docs and a cookbook for **building your own MCP server**.

### Other differentiators

- **Multi-model in one place** – Claude, GPT, Gemini, Composer 1.5, Grok, etc., with one subscription and optional BYOK.
- **Rules and AGENTS.md** – Project- and file-level instructions the AI follows automatically.
- **Composer + Agent mode** – Structured multi-file edits with diffs/checkpoints, plus an optional “agent” mode that explores and runs terminal/browser.
- **Native terminal and browser tools** – Agent can run commands and use a browser (e.g. for testing) inside Cursor.
- **Integrations** – GitHub, GitLab, Linear, Slack, etc., for PRs, issues, and notifications.

---

## How Cursor Compares to Other Tools

Short positioning vs other AI coding options:

| Tool | What it is | Best at | Main limitation vs Cursor |
|------|------------|---------|----------------------------|
| **Claude in VS Code** | Claude as an agent inside VS Code | Same “agent in editor” feel, Claude’s reasoning | No Cursor-style indexing/semantic search, no Tab, single-model, different rules/Composer UX |
| **GitHub Copilot (VS Code)** | Inline completion + chat in the IDE | Fast autocomplete, boilerplate, any editor | Small context (e.g. 4–8k tokens), weak codebase-wide understanding; chat is secondary to completion |
| **Claude CLI / Claude Code** | Terminal-first agent (e.g. `claude` in the shell) | Deep refactors, big context, terminal automation, MCP | No IDE; you’re in the shell. No Tab, no Cursor indexing; different workflow (delegate then review) |
| **Codex** (OpenAI) | Coding-optimized agent (GPT-5.3 Codex); ChatGPT + cloud agent + optional CLI | Agentic coding: features, bugs, PRs; terminal, tests, long-running tasks. SOTA on SWE-bench style evals | Access via ChatGPT (e.g. chatgpt.com/codex) or CLI—not a full IDE. Cloud agent is “delegate then review”; no Cursor-style Tab/indexing in one app. |
| **Devin** | Autonomous AI engineer (e.g. via Slack) | Hands-off tasks: plan → code → PR with minimal oversight | Very different model: you *delegate*; it runs in the background. No real-time co-editing. Much higher cost (~$500/mo). Not an IDE. |
| **Gemini for coding** | Google’s AI in IDE or elsewhere | Long context, multimodal (e.g. images), Google ecosystem | Depends on product (IDE extension vs web). Typically no Cursor-level indexing/Tab/Composer in one place. |
| **Windsurf** (Codeium) | AI-native IDE with “Cascade” agent | Agentic multi-file edits, smart autocomplete, 70+ languages, terminal integration | Closest to Cursor in “IDE + agent” feel. Different indexing/context model; no Cursor Tab; different pricing (e.g. Free / $15 Pro / $60 Ultimate). |
| **Cody** (Sourcegraph) | IDE extension (VS Code, JetBrains) with codebase AI | Code graph + @-mentions for precise context; supports Claude, GPT, custom models | Strong codebase awareness via graph, but no full IDE ownership—extension inside your editor. No Cursor-style semantic index or Tab. |
| **Replit** | Cloud IDE + AI (Agent, Complete, Chat) | Full-stack in the browser; good for learning, prototypes, and API-driven workflows | Entirely cloud-based; you don’t run the editor locally. Different workflow than desktop IDE + agent. |

**Cursor’s niche:** An **IDE-first** experience with **indexed semantic search**, **codebase-aware Tab**, **Composer + Agent**, **multi-model**, and **rules/MCP** in one product. You stay in the editor with tight feedback loops, while the AI has a structured view of the whole repo and optional external tools via MCP.

---

## Modes: When to Use Each

### 1. Chat (sidebar)

- **Use for:** Understanding code, debugging, small refactors, architecture questions, “what does this do?”, exploring the codebase.
- **Behavior:** Conversational; you apply suggested code yourself (or use apply when offered).
- **Tip:** Use **⌘⏎** to run a codebase search from chat.

### 2. Composer (⌘I)

- **Use for:** Writing and editing code: new features, multi-file changes, scaffolding, “build this” / “refactor that.” You get diffs to accept or reject.
- **Behavior:** Generates and applies edits; you review diffs and can use checkpoints to undo.
- **Layouts:** Panel (chat + editor) or full Editor tab.

**Inside Composer:**

| | Normal Composer | Agent mode (⌘. in Composer) |
|---|----------------|----------------------------|
| **When** | You know what you want; direct, predictable edits. | AI should explore the repo and do a multi-step task with less hand-holding. |
| **Behavior** | Follows your instructions and makes targeted edits. | More autonomous: gathers context (e.g. `@Recommended`), runs terminal, searches code, creates/edits files. |
| **Models** | Any Composer model. | **Claude only.** Has a tool-call limit (e.g. 25). |

### 3. Inline edit (⌘K)

- **Use for:** Single-spot edits: “add error handling here,” “rename this to X,” “extract this into a function.”
- **Behavior:** One diff in the current file; accept or reject.

### Summary

- **Explain / debug / small fix** → **Chat**
- **Change one place** → **⌘K** (inline edit)
- **Implement a feature or refactor across files** → **Composer** (⌘I), normal mode
- **“Do this whole task and figure out the steps”** → **Composer** (⌘I) + **Agent** (⌘.) with a **Claude** model

---

## Models: What’s Available and When to Use Them

Choose the model in the **model dropdown** under the input (Chat, Composer, and for ⌘K when relevant).

| Model | When to use it |
|-------|-----------------|
| **Claude 4.6 Opus** | Complex features, architecture, refactors, agent-style work. Best quality; higher cost. |
| **Claude 4.5 Sonnet** | Day-to-day work; strong quality, lower cost than Opus. |
| **Composer 1.5** | Speed-critical iteration; Cursor’s own model. |
| **GPT-5.2** | General coding; often good for UI and bug fixes. |
| **GPT-5.3 Codex** | Coding-focused OpenAI model. |
| **Gemini 3 Pro** | Very large context (e.g. 1M tokens) or multimodal (e.g. images). |
| **Gemini 3 Flash** | Lighter, cheaper option. |
| **Grok Code** | Low-cost option for simpler or high-volume edits. |
| **Auto** | Let Cursor choose the model for the current task. |

**Notes:**

- **Agent mode** (⌘. in Composer) only works with **Claude** models.
- **Max context:** Some models support a larger “Max” context (e.g. 1M tokens); use for huge codebases or long threads (slower, more usage).
- **Cost:** Opus > Sonnet > Composer 1.5; Gemini Flash and Grok are cheaper. Check **Settings → usage** and the [Cursor pricing docs](https://cursor.com/docs/account/pricing) for your plan.

---

## Pricing, Usage, and BYOK

### Your $20/month (Pro plan)

- Pro includes **$20 of included usage** per month (plus possible bonus), charged at **each model’s API rate** (or Auto’s rate if you use Auto).
- **Cursor only “decides” the model if you selected “Auto”** in the model picker; otherwise the model is the one you chose.
- When you use a model (or Auto), that usage **draws from your $20** (then overage if you have it enabled).
- View usage: **Cursor dashboard → Usage** and in-editor indicators.

### Bring your own key (BYOK)

- You can add your own API keys (OpenAI, Anthropic, Google, etc.) in **Settings → Models** (or API Keys).
- When you use a **model billed through your key**, that usage is **billed by the provider to you**, not from your Cursor $20.
- **You still need an active Cursor subscription** for full features. With only BYOK and no subscription, Composer, Apply, and Tab are not available.
- **Use BYOK if:** You already pay for a provider and want that usage to go to your account, or you want to preserve the $20 for other models.

---

## Tokens and Context: How It Works

Cursor doesn’t use slash commands like Claude’s `/clear`. Instead you manage context and tokens with built-in actions and habits.

### How tokens work

- **Token** ≈ ¾ of an English word, or ~4 characters. Models “see” and price by tokens, not words.
- **Context window** = the model’s working memory for this conversation. Most models use **200k tokens** by default (~15k lines of code); some support **Max mode** (e.g. 1M tokens).
- **What uses tokens:**
  - **Input:** Your messages, conversation history, @-mentioned files, rules, system prompt, tool outputs (e.g. terminal, linter).
  - **Output:** Everything the model generates (replies, code, tool calls).
  - **Cache:** Cursor can cache repeated context; cache read is cheaper than full input.

Output tokens usually cost more than input (generation is heavier than reading). Your **Usage** in the dashboard and in-editor shows consumption against your plan.

### “Clearing” context (Cursor’s version of /clear)

| Goal | What to do |
|------|------------|
| **Reset context but keep history visible** | **⌘⇧P** (Mac) or **Ctrl+Shift+P** → type **“Cursor: Clear Context”** → Enter. Next response won’t use the previous “active” context; history stays for reference. |
| **Start a brand‑new conversation** | **New Chat / New Composer:** Use the **+** or “New chat” in the Chat panel, or **⌘N** for a new Composer. Fresh thread = no prior messages in context. |
| **Remove specific files from context** | Hover over the **@ chips** (file pills) above the input and click **✕** on the ones you don’t need. Fewer @ files = fewer tokens. |
| **Chat with minimal auto context** | **Settings → Features → Chat** → enable **“Default No Context”** so the model doesn’t auto-include open file etc. unless you @ something. |

So: there’s no `/clear` command; you **clear context** (command palette), **start a new chat/composer** (new conversation), or **prune @ mentions** (remove file chips).

### Staying efficient with tokens

- **Prune, don’t nuke:** Remove only the heavy @ file references you no longer need instead of deleting the whole chat when you hit limits.
- **Avoid `@Codebase`** unless necessary; **@ specific files or folders** to send less and be more precise.
- **Start a new Composer/Chat** when switching to a different feature or task so the thread doesn’t grow unbounded.
- **Use the context indicator** in the UI (e.g. “Context: 0/200k”) to see how full the window is.
- **Rules and AGENTS.md** add tokens every turn; keep them concise so they don’t dominate the budget.

---

## Project Context: AGENTS.md and Rules

### AGENTS.md (project root)

A **single file** that tells the AI how to work in **this** codebase. Put at the repo root.

**Include:**

- Tech stack (languages, frameworks, runtimes).
- Project structure (where routes, lib, API, etc. live).
- Conventions (naming, file layout, patterns).
- Do’s and don’ts (e.g. “Use our `Button` from `@/components/ui`”).
- How to run/test (commands for dev, tests, lint, migrations).
- Any other “must know” details (env, scripts, defaults).

Keep it **short and scannable** (bullets, clear sections).

### .cursor/rules/ (optional)

- **Multiple** small rule files (`.mdc`) with YAML frontmatter.
- Use for **“when editing these files, do this”** (e.g. by file type via `globs` or `alwaysApply`).
- **AGENTS.md** = project-wide “what this project is”; **rules** = targeted, often file-specific behavior.

---

## Quick Reference

| Goal | Use |
|------|-----|
| Explain / debug / small fix | **Chat** |
| Change one spot in a file | **⌘K** (inline edit) |
| Implement feature / refactor across files | **Composer** (⌘I), normal mode |
| “Do this whole task” with autonomy | **Composer** (⌘I) + **Agent** (⌘.) with **Claude** |
| Give the AI project context | **AGENTS.md** + optional **.cursor/rules/** |
| Point AI at specific files/folders | **@** mentions in Chat or Composer |
| Control cost / use your own credits | Pick a specific model (or Auto); optionally **BYOK** for that provider |
| Clear context / save tokens | **⌘⇧P** → “Cursor: Clear Context”; or new chat / **⌘N**; or remove **@** file chips |
| See token usage | In-editor context indicator (e.g. 0/200k); **Cursor dashboard → Usage** |

---

*Last updated for Cursor as of February 2025. Check [Cursor Docs](https://cursor.com/docs) for current behavior and pricing.*
