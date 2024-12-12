# MCP

MCP is a command-line tool and local UI for discovering, installing and managing [Model Context Protocol](https://modelcontextprotocol.io/) servers.

> [!NOTE]
> This README describes an aspirational vision. Almost none of this is ready but this is intended to set the direction we're going.

The `mcp` tool acts as an MCP Server for Clients like [Claude](https://claude.ai/) and [Zed](https://zed.dev/). It is unique in that it doesn't provide any capabilities on its own. Instead it acts as a sort of broker between Clients and any number of installed Servers. The `mcp` tool facilitates that by helping you discover and install MCP Servers. It talks to a public registry of MCP Servers and is able to download, configure and run these in a secure and frictionless way.

In addition to acting as a broker, `mcp`:

1. Allows Servers to request OAuth2 credentials and manages the flow for obtaining, storing and refreshing these credentials.
1. Keeps an audit log of all operations going through it and provides a UI for reviewing audit logs across different sessions and Clients.
1. Reduces host system dependencies to `docker`. By leveraging `docker`, `mcp` allows you to run Servers without mutating your host with any necessary tool-chains or dependencies. In addition, `mcp` leverages the containerization and isolation technology in `docker` to isolate your host system from malicious or buggy MCP Servers.

# CLI Interface

## mcp install <claude|zed|...>

Install `mcp run <protocol>` as an MCP Server for the supplied Client.

## mcp registry search <query>

Search the public MCP server Registry for servers matching the query. A rich description will be shown for matching servers.

## mcp server install <server[@<version>]>

Install an MCP server from the public server Registry. This will start a flow that captures any required configuration for the MCP Server, persist it locally and then start it.

## mcp server uninstall <server>

Uninstall an MCP Server that was previously installed. Running clients will be notified such that they reload resources, tools, etc.

## mcp run stdio

This is the entrypoint used by Clients that speak the `stdio` protocol. It will run `mcp` as an MCP Server that acts as a broker for all installed MCP Servers.
