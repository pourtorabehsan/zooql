import ConnectionDetails from "./connection_details";
import QueryPanel from "./query_panel";

export default function Main() {
  return (
    <div className="flex flex-col flex-1">
      <ConnectionDetails />
      <QueryPanel />
    </div>
  );
}
