# [FEATURE] Implement Search Functionality with `/` Key

## 📝 Feature Description
Add search functionality to the TUI application that allows users to quickly find tasks by pressing the `/` key, similar to Vim's search functionality.

## 🎯 Problem to Solve
Currently, users need to manually scroll through all tasks to find specific items, which becomes inefficient when dealing with large todo files:
- No quick way to locate specific tasks
- Time-consuming to find tasks containing specific keywords
- Difficult to search across different projects, contexts, or due dates
- Users coming from Vim expect familiar search patterns

## 💡 Proposed Solution
Implement a search mode that activates when the user presses `/`:

### 1. Search Mode Activation
- Press `/` to enter search mode
- Display search input field at the bottom of the screen
- Show search prompt (e.g., "Search: ")

### 2. Search Input Handling
- Accept text input for search query
- Support real-time filtering as user types
- Handle special characters and regex patterns
- Clear existing filters when entering search mode

### 3. Search Functionality
- **Text Search**: Find tasks containing specific text
- **Project Search**: Search for specific projects (e.g., `+project`)
- **Context Search**: Search for specific contexts (e.g., `@context`)
- **Priority Search**: Search for specific priorities (e.g., `(A)`)
- **Date Search**: Search for specific dates or date ranges

### 4. Search Navigation
- `Enter`: Apply search filter and exit search mode
- `Esc`: Cancel search and return to normal mode
- Arrow keys or Vim keys: Navigate through filtered results
- `n`/`N`: Jump to next/previous search result (after applying search)

### 5. Visual Feedback
- Highlight search terms in task list
- Show number of matches found
- Display "No matches found" when appropriate
- Maintain highlight until new search or clear action

## 🔄 Alternatives Considered
1. **Filter-only approach**: Use existing filter pane (less intuitive for quick searches)
2. **Command palette**: Use Ctrl+P style search (less familiar to terminal users)
3. **Multiple search keys**: Use different keys for different search types (more complex)

## ✅ Acceptance Criteria
- [ ] `/` key activates search mode from normal mode
- [ ] Search input field appears at bottom of screen
- [ ] Real-time filtering works as user types
- [ ] Search works for task text content
- [ ] Search works for projects (`+project`)
- [ ] Search works for contexts (`@context`)
- [ ] Search works for priorities (`(A)`, `(B)`, etc.)
- [ ] `Enter` applies search and exits search mode
- [ ] `Esc` cancels search and returns to normal mode
- [ ] Search terms are highlighted in results
- [ ] Match count is displayed
- [ ] `n`/`N` navigation works for search results
- [ ] Search is case-insensitive by default
- [ ] Empty search shows all tasks
- [ ] Search state is cleared appropriately

## 📱 Additional Context
### Current Key Bindings (from README):
```
j/k     - Navigate lists
Tab     - Switch between filter and task panes
Enter   - Apply filter / Complete task
a       - Add new task
e       - Edit task
d       - Delete task
p       - Cycle priority
r       - Restore deleted/completed task
q       - Quit
```

### Search Examples:
- `/buy` - Find tasks containing "buy"
- `/+groceries` - Find tasks in groceries project
- `/@store` - Find tasks with store context
- `/(A)` - Find high priority tasks
- `/due:2025` - Find tasks due in 2025

### UI Mockup:
```
┌─ Filters ─────────────┬─ Tasks ─────────────────────────┐
│ Projects              │ (A) Call Mom @phone +family     │
│ + family (2)          │ Buy milk @store +groceries      │
│ + groceries (1)       │ Clean garage @home +chores      │
│                       │                                 │
│ Contexts              │                                 │
│ @ phone (1)           │                                 │
│ @ store (1)           │                                 │
│ @ home (1)            │                                 │
└───────────────────────┴─────────────────────────────────┘
Search: buy                                    1 match found
```

## 🏷️ Priority
- [x] Medium (should be addressed soon)

## 📋 Implementation Tasks
- [ ] Add search mode state to TUI model
- [ ] Implement `/` key handler to activate search mode
- [ ] Create search input component
- [ ] Implement real-time search filtering
- [ ] Add search highlighting functionality
- [ ] Implement `Enter`/`Esc` handlers for search mode
- [ ] Add `n`/`N` navigation for search results
- [ ] Add match counter display
- [ ] Update help/documentation with search functionality
- [ ] Add unit tests for search functionality
- [ ] Test search with various todo.txt formats

## 🔗 Related Issues
- Depends on current filtering system
- May enhance existing filter functionality
- Consider integration with existing key bindings 