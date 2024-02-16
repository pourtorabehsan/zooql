import axios from "axios";
import useSWR from "swr";

export default function ConnectionDetails() {
  const fetcher = (...args) => axios.get(...args).then((res) => res.data);
  const { data, error, isLoading } = useSWR("/backend/connection", fetcher);

  return (
    <div className="p-8 border-b flex flex-col gap-2">
      <div className="flex gap-2">
        <span>Connection:</span>
        <span className="text-success font-code">
          {isLoading ? "..." : error ? `Error: ${error}` : data.zookeepers}
        </span>
      </div>
      <div className="flex gap-2">
        <span>Base Path:</span>
        <span className="text-success font-code">
          {isLoading ? "..." : error ? `Error: ${error}` : data.basePath}
        </span>
      </div>
    </div>
  );
}
