import { useState } from "react";

export default function QueryInput({ onSubmit, history }) {
  const [query, setQuery] = useState("");
  const [historyIndex, setHistoryIndex] = useState(history.length);

  const handleKeyUp = async (e) => {
    if (e.key === "Enter") {
      const trimmedQuery = query.trim();
      if (trimmedQuery && onSubmit) {
        await onSubmit(trimmedQuery);
      }
    }
  };

  const handleKeyDown = (e) => {
    if (e.key === "ArrowUp") {
      if (history.length > 0) {
        const nextIndex = Math.max(0, historyIndex - 1);
        const nextQuery = history[nextIndex];
        setHistoryIndex(nextIndex);
        if (nextQuery) {
          setQuery(nextQuery);
        }
      }
    }

    if (e.key === "ArrowDown") {
      if (history.length > 0) {
        const nextIndex = Math.min(history.length, historyIndex + 1);
        const nextQuery = history[nextIndex];
        setHistoryIndex(nextIndex);
        if (nextQuery) {
          setQuery(nextQuery);
        } else {
          setQuery("");
        }
      }
    }
  };

  return (
    <input
      className="flex-1 bg-transparent border-transparent focus:outline-none font-code text-lg p-2"
      value={query}
      onChange={(e) => setQuery(e.target.value)}
      onKeyUp={handleKeyUp}
      onKeyDown={handleKeyDown}
    ></input>
  );
}
