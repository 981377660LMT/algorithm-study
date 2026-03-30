#!/bin/bash
INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name')
TOOL_INPUT=$(echo "$INPUT" | jq -r '.tool_input')

if [ "$TOOL_NAME" = "runTerminalCommand" ]; then
  COMMAND=$(echo "$TOOL_INPUT" | jq -r '.command // empty')

  if echo "$COMMAND" | grep -qE '(rm\s+-rf|DROP\s+TABLE|DELETE\s+FROM)'; then
    echo '{"hookSpecificOutput":{"permissionDecision":"deny","permissionDecisionReason":"Destructive command blocked by security policy"}}'
    exit 0
  fi
fi

echo '{"continue":true}'
