import Image from "next/image";

export default function Header() {
  return (
    <header className="flex items-center gap-2 border-b p-8">
      <Image src="/logo.png" alt="ZooQL Logo" width={32} height={32} />
      <h1 className="text-3xl font-bold font-code">ZooQL</h1>
    </header>
  );
}
