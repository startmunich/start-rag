import Dagre from '@dagrejs/dagre';
import { Node, Edge } from 'reactflow';
import { CrawledPage } from "@/components/graph/graph";

export const maxPages = (pages: CrawledPage[], max: number): CrawledPage[] => {
    return pages.slice(0, max);
}

export const parsePagesToGraph = (pages: CrawledPage[]): { nodes: Node[], edges: Edge[] } => {
    const nodes: Node[] = [];
    const edges: Edge[] = [];

    for (const page of pages) {
        nodes.push({
            id: page.page_id,
            type: "crawledPage",
            data: page,
            position: { x: Math.random() * 1000, y: Math.random() * 1000 },
        });

        for (const childPage of page.child_pages.split(';')) {
            if (pages.find((p) => p.page_id === childPage) === undefined) continue;
            edges.push({
                id: `${page.page_id}-${childPage}`,
                type: 'smoothstep',
                source: page.page_id,
                target: childPage,
            });
        }
    }

    return {
        nodes: nodes,
        edges: edges,
    };
}

const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));

export const getLayoutedElements = (nodes: Node[], edges: Edge[]): { nodes: Node[], edges: Edge[] } => {
  g.setGraph({ rankdir: "TB", ranksep: 100, nodesep: 120});

  edges.forEach((edge) => g.setEdge(edge.source, edge.target));
  nodes.forEach((node) => g.setNode(node.id, {}));

  Dagre.layout(g, { weight: 1 });

  return {
    nodes: nodes.map((node) => {
      const { x, y } = g.node(node.id);

      return { ...node, position: { x, y } };
    }),
    edges,
  };
};
