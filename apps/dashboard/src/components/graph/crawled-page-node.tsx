import { Handle, NodeProps, Position } from 'reactflow';
import { CrawledPage } from './graph';
import { NotionLogoIcon, OpenInNewWindowIcon } from '@radix-ui/react-icons';
import Link from 'next/link';
 
const handleStyle = { left: 10 };
 
export function CrawledPageNode({ data }: NodeProps) {
  const { page_id, url } = data as CrawledPage;
  return (
    <div className="w-28 h-8 bg-[#1e2024] border-[#383a3e] border-2 border-solid rounded-lg text-white">
      <Handle type="target" position={Position.Top} />
      <div className="flex justify-start px-2 items-center h-full">
        <NotionLogoIcon />
        <span className="px-1 pt-0.5 text-ellipsis text-xs flex-grow">{page_id.substring(0,5)}</span>
        <Link target="_blank" className="rounded-full bg-blue-300/50 py-0.5 px-2" href={url}>
          <OpenInNewWindowIcon />
        </Link>
      </div>
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
}