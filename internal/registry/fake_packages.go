package registry

import "encoding/json"

type mcpGetPackage struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Vendor      string `json:"vendor"`
	SourceURL   string `json:"sourceUrl"`
	Homepage    string `json:"homepage"`
	Licence     string `json:"licence"`
	Runtime     string `json:"runtime"`
}

var rawData = []byte(`[
  {
    "name": "@modelcontextprotocol/server-brave-search",
    "description": "MCP server for Brave Search API integration",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/brave-search",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-everything",
    "description": "MCP server that exercises all the features of the MCP protocol",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/everything",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-filesystem",
    "description": "MCP server for filesystem access",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/filesystem",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-gdrive",
    "description": "MCP server for interacting with Google Drive",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/gdrive",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-github",
    "description": "MCP server for using the GitHub API",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/github",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-gitlab",
    "description": "MCP server for using the GitLab API",
    "vendor": "GitLab, PBC (https://gitlab.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/gitlab",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-google-maps",
    "description": "MCP server for using the Google Maps API",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/google-maps",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-memory",
    "description": "MCP server for enabling memory for Claude through a knowledge graph",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/memory",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-postgres",
    "description": "MCP server for interacting with PostgreSQL databases",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/postgres",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-puppeteer",
    "description": "MCP server for browser automation using Puppeteer",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/puppeteer",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-slack",
    "description": "MCP server for interacting with Slack",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/slack",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@cloudflare/mcp-server-cloudflare",
    "description": "MCP server for interacting with Cloudflare API",
    "vendor": "Cloudflare, Inc. (https://cloudflare.com)",
    "sourceUrl": "https://github.com/cloudflare/mcp-server-cloudflare",
    "homepage": "https://github.com/cloudflare/mcp-server-cloudflare",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@raygun.io/mcp-server-raygun",
    "description": "MCP server for interacting with Raygun's API for crash reporting and real user monitoring metrics",
    "vendor": "Raygun (https://raygun.com)",
    "sourceUrl": "https://github.com/MindscapeHQ/mcp-server-raygun",
    "homepage": "https://raygun.com",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@kimtaeyoon83/mcp-server-youtube-transcript",
    "description": "This is an MCP server that allows you to directly download transcripts of YouTube videos.",
    "vendor": "Freddie (https://github.com/kimtaeyoon83)",
    "sourceUrl": "https://github.com/kimtaeyoon83/mcp-server-youtube-transcript",
    "homepage": "https://github.com/kimtaeyoon83/mcp-server-youtube-transcript",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@kagi/mcp-server-kagi",
    "description": "MCP server for Kagi search API integration",
    "vendor": "ac3xx (https://github.com/ac3xx)",
    "sourceUrl": "https://github.com/ac3xx/mcp-servers-kagi",
    "homepage": "https://github.com/ac3xx/mcp-servers-kagi",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@exa/mcp-server",
    "description": "MCP server for Exa AI Search API integration",
    "vendor": "Exa Labs (https://exa.ai)",
    "sourceUrl": "https://github.com/exa-labs/exa-mcp-server",
    "homepage": "https://exa.ai",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@search1api/mcp-server",
    "description": "MCP server for Search1API integration",
    "vendor": "fatwang2 (https://github.com/fatwang2)",
    "sourceUrl": "https://github.com/fatwang2/search1api-mcp",
    "homepage": "https://github.com/fatwang2/search1api-mcp",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@calclavia/mcp-obsidian",
    "description": "MCP server for reading and searching Markdown notes (like Obsidian vaults)",
    "vendor": "Calclavia (https://github.com/calclavia)",
    "sourceUrl": "https://github.com/calclavia/mcp-obsidian",
    "homepage": "https://github.com/calclavia/mcp-obsidian",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@anaisbetts/mcp-youtube",
    "description": "MCP server for fetching YouTube subtitles",
    "vendor": "Anaïs Betts (https://github.com/anaisbetts)",
    "sourceUrl": "https://github.com/anaisbetts/mcp-youtube",
    "homepage": "https://github.com/anaisbetts/mcp-youtube",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-everart",
    "description": "MCP server for EverArt API integration",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/everart",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-sequential-thinking",
    "description": "MCP server for sequential thinking and problem solving",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/sequentialthinking",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "mcp-server-fetch",
    "description": "A Model Context Protocol server providing tools to fetch and convert web content for usage by LLMs",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/fetch",
    "homepage": "https://github.com/modelcontextprotocol/servers",
    "license": "MIT",
    "runtime": "python"
  },
  {
    "name": "mcp-server-perplexity",
    "description": "MCP Server for the Perplexity API",
    "vendor": "tanigami",
    "sourceUrl": "https://github.com/tanigami/mcp-server-perplexity",
    "homepage": "https://github.com/tanigami/mcp-server-perplexity",
    "license": "MIT",
    "runtime": "python"
  },
  {
    "name": "mcp-server-git",
    "description": "A Model Context Protocol server providing tools to read, search, and manipulate Git repositories programmatically via LLMs",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/git",
    "homepage": "https://github.com/modelcontextprotocol/servers",
    "license": "MIT",
    "runtime": "python"
  },
  {
    "name": "mcp-server-sentry",
    "description": "MCP server for retrieving issues from sentry.io",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/sentry",
    "homepage": "https://github.com/modelcontextprotocol/servers",
    "license": "MIT",
    "runtime": "python"
  },
  {
    "name": "mcp-server-sqlite",
    "description": "A simple SQLite MCP server",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/sqlite",
    "homepage": "https://github.com/modelcontextprotocol/servers",
    "license": "MIT",
    "runtime": "python"
  },
  {
    "name": "mcp-server-time",
    "description": "A Model Context Protocol server providing tools for time queries and timezone conversions for LLMs",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/time",
    "homepage": "https://github.com/modelcontextprotocol/servers",
    "license": "MIT",
    "runtime": "python"
  },
  {
    "name": "mcp-tinybird",
    "description": "A Model Context Protocol server that lets you interact with a Tinybird Workspace from any MCP client.",
    "vendor": "Tinybird (https://tinybird.co)",
    "sourceUrl": "https://github.com/tinybirdco/mcp-tinybird/tree/main/src/mcp-tinybird",
    "homepage": "https://github.com/tinybirdco/mcp-tinybird",
    "license": "Apache 2.0",
    "runtime": "python"
  },
  {
    "name": "@automatalabs/mcp-server-playwright",
    "description": "MCP server for browser automation using Playwright",
    "vendor": "Automata Labs (https://automatalabs.io)",
    "sourceUrl": "https://github.com/Automata-Labs-team/MCP-Server-Playwright/tree/main",
    "homepage": "https://github.com/Automata-Labs-team/MCP-Server-Playwright",
    "runtime": "node",
    "license": "MIT"
  },
  {
    "name": "@mcp-get-community/server-llm-txt",
    "description": "MCP server that extracts and serves context from llm.txt files, enabling AI models to understand file structure, dependencies, and code relationships in development environments",
    "vendor": "Michael Latman (https://michaellatman.com)",
    "sourceUrl": "https://github.com/mcp-get/community-servers/blob/main/src/server-llm-txt",
    "homepage": "https://github.com/mcp-get/community-servers#readme",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@executeautomation/playwright-mcp-server",
    "description": "A Model Context Protocol server for Playwright for Browser Automation and Web Scraping.",
    "vendor": "ExecuteAutomation, Ltd (https://executeautomation.com)",
    "sourceUrl": "https://github.com/executeautomation/mcp-playwright/tree/main/src",
    "homepage": "https://github.com/executeautomation/mcp-playwright",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@mcp-get-community/server-curl",
    "description": "MCP server for making HTTP requests using a curl-like interface",
    "vendor": "Michael Latman <https://michaellatman.com>",
    "sourceUrl": "https://github.com/mcp-get/community-servers/blob/main/src/server-curl",
    "homepage": "https://github.com/mcp-get-community/server-curl#readme",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@mcp-get-community/server-macos",
    "description": "MCP server for macOS system operations",
    "vendor": "Michael Latman <https://michaellatman.com>",
    "sourceUrl": "https://github.com/mcp-get/community-servers/blob/main/src/server-macos",
    "homepage": "https://github.com/mcp-get-community/server-macos#readme",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@modelcontextprotocol/server-aws-kb-retrieval",
    "description": "MCP server for AWS Knowledge Base retrieval using Bedrock Agent Runtime",
    "vendor": "Anthropic, PBC (https://anthropic.com)",
    "sourceUrl": "https://github.com/modelcontextprotocol/servers/blob/main/src/aws-kb-retrieval-server",
    "homepage": "https://modelcontextprotocol.io",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "docker-mcp",
    "description": "A powerful Model Context Protocol (MCP) server for Docker operations, enabling seamless container and compose stack management through Claude AI",
    "vendor": "QuantGeekDev & md-archive",
    "sourceUrl": "https://github.com/QuantGeekDev/docker-mcp",
    "homepage": "https://github.com/QuantGeekDev/docker-mcp",
    "license": "MIT",
    "runtime": "python"
  },
  {
    "name": "mcp-mongo-server",
    "description": "A Model Context Protocol Server for MongoDB",
    "vendor": "Muhammed Kılıç <kiliczsh>",
    "sourceUrl": "https://github.com/kiliczsh/mcp-mongo-server",
    "homepage": "https://github.com/kiliczsh/mcp-mongo-server",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@llmindset/mcp-hfspace",
    "description": "MCP Server for using HuggingFace Spaces. Seamlessly use the latest Open Source Image, Audio and Text Models from within Claude Deskop.",
    "vendor": "llmindset.co.uk",
    "sourceUrl": "https://github.com/evalstate/mcp-hfspace/",
    "homepage": "https://llmindset.co.uk/resources/hfspace-connector/",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@llmindset/mcp-miro",
    "description": "A Model Context Protocol server to connect to the MIRO Whiteboard Application",
    "vendor": "llmindset.co.uk",
    "sourceUrl": "https://github.com/evalstate/mcp-miro",
    "homepage": "https://github.com/evalstate/mcp-miro#readme",
    "license": "Apache-2.0",
    "runtime": "node"
  },
  {
    "name": "@strowk/mcp-k8s",
    "description": "MCP server connecting to Kubernetes",
    "vendor": "Timur Sultanaev (https://str4.io/about-me)",
    "sourceUrl": "https://github.com/strowk/mcp-k8s-go",
    "homepage": "https://github.com/strowk/mcp-k8s-go",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "mcp-shell",
    "description": "An MCP server for your shell",
    "vendor": "High Dimensional Research (https://hdr.is)",
    "sourceUrl": "https://github.com/hdresearch/mcp-shell",
    "homepage": "https://github.com/hdresearch/mcp-shell",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "@benborla29/mcp-server-mysql",
    "description": "An MCP server for interacting with MySQL databases",
    "vendor": "Ben Borla (https://benborla.dev)",
    "sourceUrl": "https://github.com/benborla/mcp-server-mysql",
    "homepage": "https://github.com/benborla/mcp-server-mysql",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "mcp-server-rememberizer",
    "description": "An MCP server for interacting with Rememberizer's document and knowledge management API. This server enables Large Language Models to search, retrieve, and manage documents and integrations through Rememberizer.",
    "vendor": "Rememberizer®",
    "sourceUrl": "https://github.com/skydeckai/mcp-server-rememberizer",
    "homepage": "https://rememberizer.ai/",
    "license": "MIT",
    "runtime": "python"
  },
  {
    "name": "@enescinar/twitter-mcp",
    "description": "This MCP server allows Clients to interact with Twitter, enabling posting tweets and searching Twitter.",
    "vendor": "Enes Çınar",
    "sourceUrl": "https://github.com/EnesCinr/twitter-mcp",
    "homepage": "https://github.com/EnesCinr/twitter-mcp",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "mcp-server-commands",
    "description": "MCP server enabling LLMs to execute shell commands and run scripts through various interpreters with built-in safety controls.",
    "vendor": "g0t4 (https://github.com/g0t4)",
    "sourceUrl": "https://github.com/g0t4/mcp-server-commands",
    "homepage": "https://github.com/g0t4/mcp-server-commands",
    "license": "MIT",
    "runtime": "node"
  },
  {
    "name": "mcp-server-kubernetes",
    "description": "MCP server for managing Kubernetes clusters, enabling LLMs to interact with and control Kubernetes resources.",
    "vendor": "Flux159 (https://github.com/Flux159)",
    "sourceUrl": "https://github.com/Flux159/mcp-server-kubernetes",
    "homepage": "https://github.com/Flux159/mcp-server-kubernetes",
    "license": "MIT",
    "runtime": "node"
  }
]`)

func loadFakePackages() ([]mcpGetPackage, error) {
	var pkgs []mcpGetPackage

	if err := json.Unmarshal(rawData, &pkgs); err != nil {
		return nil, err
	}

	return pkgs, nil
}
