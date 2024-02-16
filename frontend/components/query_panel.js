import { useState } from "react";
import QueryInput from "./query_input";
import QueryResult from "./query_result";
import axios from "axios";
import { Alert, AlertDescription } from "./ui/alert";
import { Button } from "./ui/button";
import { DownloadIcon } from "lucide-react";

export default function QueryPanel() {
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);
  const [history, setHistory] = useState([]);

  const handleQuerySubmit = async (query) => {
    try {
      const res = await axios.post("/backend/query", { query });

      setResult(res.data);
      setHistory([...history, query]);
      setError(null);
    } catch (error) {
      console.log(error);
      setResult(null);
      setError(error.response.data || error);
    }
  };

  const handleExportCsv = () => {
    let csv = result.columns.join(",");
    csv = csv + "\n" + result.rows.map((row) => row.join(",")).join("\n");

    const blob = new Blob([csv], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.setAttribute("download", "export.csv");
    a.click();
  };

  return (
    <div className="flex flex-col gap-4 p-8">
      <div className="flex items-center">
        <span className="text-info text-lg font-code">SQL {">"}</span>
        <QueryInput onSubmit={handleQuerySubmit} history={history} />
      </div>
      {result && (
        <div className="flex flex-col gap-2">
          <div className="border rounded-md">
            <div className="bg-muted text-muted-foreground px-4 py-2 flex justify-between items-center">
              <div>
                {result.rows.length} rows in ({result.elapsed})
              </div>
              <div>
                <Button variant="ghost" onClick={handleExportCsv}>
                  <DownloadIcon className="w-4 h-4 mr-2" /> Export CSV
                </Button>
              </div>
            </div>
            <QueryResult result={result} />
          </div>
        </div>
      )}
      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
    </div>
  );
}
