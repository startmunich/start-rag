import CrawledPagesTable from "@/components/crawled-pages-table/crawled-pages-table";

export default function Home() {
  return (
    <main className="h-[100vh]">
      <div className="flex items-center justify-center mt-7">
        <h1 className="text-3xl font-bold text-center text-gray-900">
          Workspace
        </h1>
      </div>
      <div className="flex justify-center gap-4 mx-7 mt-7">
        <div className="w-full max-w-[950px] relative">
          <CrawledPagesTable />
        </div>
      </div>
    </main>
  );
}
