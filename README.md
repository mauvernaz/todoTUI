# 📝 To-Do TUI

A minimalist, clean Terminal User Interface for managing tasks, built with Go and Bubble Tea.

## 🚀 How to Start

### Option 1: Run the pre-compiled executable (Recommended)
You can run the app immediately without needing Go in your PATH:
```powershell
.\todotui.exe
```

### Option 2: Run from source
If you want to run via Go:
```powershell
& "C:\Program Files\Go\bin\go.exe" run .
```

---

## 📖 Commands & Help

Once inside the app, you can press **`?`** or **`h`** to see the full help menu.

### Navigation
- `↑` or `k`: Move selection up
- `↓` or `j`: Move selection down

### Task Management
- `n` or `a`: **Add** a new task (enters Input Mode)
- `x`, `d`, or `Backspace`: **Delete** selected task
- `Enter`: **Save** task (when in Input Mode)

### Application
- `?` or `h`: Toggle **Help** view
- `q` or `Esc`: **Quit** or return to list
- `Ctrl+C`: Force quit

---

## 🛠️ Development
- **Language**: Go
- **Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Input**: [Bubbles](https://github.com/charmbracelet/bubbles)
