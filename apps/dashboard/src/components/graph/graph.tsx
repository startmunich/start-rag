'use client'
import { parsePagesToGraph, getLayoutedElements, maxPages } from '@/lib/graph_utils';
import { useEffect, useState } from 'react';
import ReactFlow, { Node, Edge, Controls, Background, useNodesState, useEdgesState, NodeTypes, MiniMap } from 'reactflow';
import 'reactflow/dist/style.css';
import { CrawledPageNode } from './crawled-page-node';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';

export interface CrawledPage {
  child_pages: string;
  page_id:     string;
  url:         string;
}

function Graph() {
  const [max, setMax] = useState<number>(50);
  const [crawledPages, setCrawledPages] = useState<CrawledPage[]>([]);

  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);

  const nodeTypes: NodeTypes = { crawledPage: CrawledPageNode }

  const updateCrawledPages = async () => {
    const res = await fetch(`${process.env.NEXT_PUBLIC_CRAWLER_API_BASE_PATH!}/pages`);
    const data = await res.json();
    setCrawledPages(data.pages);
  };

  useEffect(() => {
    updateCrawledPages();
  }, []);

  useEffect(() => {
    const { nodes, edges } = parsePagesToGraph(maxPages(crawledPages, max));
    const layouted = getLayoutedElements(nodes, edges);
    setNodes(layouted.nodes);
    setEdges(layouted.edges);
  }, [crawledPages, setNodes, setEdges, max]);
  
  return (
    <div className="relative h-full">
      <ReactFlow nodeTypes={nodeTypes} nodes={nodes} edges={edges}>
        <Background />
        <Controls />
        <MiniMap />
      </ReactFlow>
      <div className="absolute top-0 right-0 p-4">
          <Select value={max.toString()} onValueChange={(v) => setMax(Number.parseInt(v))}>
            <SelectTrigger style={{ pointerEvents: "all" }} className="w-[140px] bg-white">
              <SelectValue placeholder="Max Nodes" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="10">10 Nodes</SelectItem>
              <SelectItem value="50">50 Nodes</SelectItem>
              <SelectItem value="100">100 Nodes</SelectItem>
              <SelectItem value="500">500 Nodes</SelectItem>
              <SelectItem value="1000">1000 Nodes</SelectItem>
              <SelectItem value="1500">1500 Nodes</SelectItem>
              <SelectItem value="-1">All Nodes</SelectItem>
            </SelectContent>
          </Select>
        </div>
    </div>
  );
}

export default Graph;