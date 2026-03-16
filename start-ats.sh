#!/usr/bin/env bash
# =============================================================================
# start-ats.sh — ATS Multi-Agent tmux Layout
# Usage: ./start-ats.sh [project-path]
#
# Layout:
#   Pane 0 (left, large)   : Orchestrator — INTERACTIVE
#   Pane 1 (top-right)     : PO session monitor
#   Pane 2 (mid-right)     : SA / PM session monitor
#   Pane 3 (bottom-right)  : DEV / QA session monitor
# =============================================================================

SESSION="ats-sdlc"
PROJECT_DIR="${1:-$(pwd)}"

# -- sanity check -------------------------------------------------------------
if ! command -v tmux &>/dev/null; then
  echo "❌  tmux not found. Install with: sudo apt install tmux  (or brew install tmux)"
  exit 1
fi

if ! command -v opencode &>/dev/null; then
  echo "❌  opencode not found. Install with: curl -fsSL https://opencode.ai/install | bash"
  exit 1
fi

# -- kill existing session if any ---------------------------------------------
tmux kill-session -t "$SESSION" 2>/dev/null

# -- create session -----------------------------------------------------------
tmux new-session -d -s "$SESSION" -c "$PROJECT_DIR" -x 220 -y 50

# Rename first window
tmux rename-window -t "$SESSION:0" "ATS-SDLC"

# -- layout: left = orchestrator (60%), right = monitors (40%) ---------------
tmux split-window -h -t "$SESSION:0" -p 40

# Split right pane into 3 vertical stacks
tmux split-window -v -t "$SESSION:0.1" -p 67
tmux split-window -v -t "$SESSION:0.2" -p 50

# -- pane labels --------------------------------------------------------------
# Pane 0 — Orchestrator (interactive)
tmux send-keys -t "$SESSION:0.0" "clear && echo ''" Enter
tmux send-keys -t "$SESSION:0.0" "echo '╔══════════════════════════════════════╗'" Enter
tmux send-keys -t "$SESSION:0.0" "echo '║   🎯  ORCHESTRATOR  (interactive)   ║'" Enter
tmux send-keys -t "$SESSION:0.0" "echo '╚══════════════════════════════════════╝'" Enter
tmux send-keys -t "$SESSION:0.0" "echo 'Start: opencode (then select @sdlc-ats-orchestrator)'" Enter

# Pane 1 — PO monitor
tmux send-keys -t "$SESSION:0.1" "clear && echo ''" Enter
tmux send-keys -t "$SESSION:0.1" "echo '┌─────────────────────────────────────┐'" Enter
tmux send-keys -t "$SESSION:0.1" "echo '│  👤  PO  (monitor)                  │'" Enter
tmux send-keys -t "$SESSION:0.1" "echo '└─────────────────────────────────────┘'" Enter
tmux send-keys -t "$SESSION:0.1" "echo 'Attach: opencode attach --session ats-po-p1-<slug>'" Enter

# Pane 2 — SA / PM monitor
tmux send-keys -t "$SESSION:0.2" "clear && echo ''" Enter
tmux send-keys -t "$SESSION:0.2" "echo '┌─────────────────────────────────────┐'" Enter
tmux send-keys -t "$SESSION:0.2" "echo '│  🏗️   SA / PM  (monitor)             │'" Enter
tmux send-keys -t "$SESSION:0.2" "echo '└─────────────────────────────────────┘'" Enter
tmux send-keys -t "$SESSION:0.2" "echo 'Attach: opencode attach --session ats-sa-p2-<slug>'" Enter

# Pane 3 — DEV / QA monitor
tmux send-keys -t "$SESSION:0.3" "clear && echo ''" Enter
tmux send-keys -t "$SESSION:0.3" "echo '┌─────────────────────────────────────┐'" Enter
tmux send-keys -t "$SESSION:0.3" "echo '│  🔧  DEV / QA  (monitor)            │'" Enter
tmux send-keys -t "$SESSION:0.3" "echo '└─────────────────────────────────────┘'" Enter
tmux send-keys -t "$SESSION:0.3" "echo 'Attach: opencode attach --session ats-dev-p4-<id>'" Enter

# -- focus orchestrator pane --------------------------------------------------
tmux select-pane -t "$SESSION:0.0"

# -- attach -------------------------------------------------------------------
echo ""
echo "✅  tmux session '$SESSION' created."
echo ""
echo "Pane layout:"
echo "  [0] LEFT       — Orchestrator (interactive)"
echo "  [1] TOP-RIGHT  — PO monitor"
echo "  [2] MID-RIGHT  — SA / PM monitor"
echo "  [3] BOT-RIGHT  — DEV / QA monitor"
echo ""
echo "Attaching now..."
echo ""

tmux attach-session -t "$SESSION"
