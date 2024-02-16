import Header from "./header";

export default function Layout({ children }) {
  return (
    <div className="flex flex-col min-h-screen bg-background font-sans antialiased font-medium">
      <Header />
      <main className="flex flex-col flex-1">{children}</main>
    </div>
  );
}
