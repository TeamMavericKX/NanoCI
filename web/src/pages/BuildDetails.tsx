import { useEffect, useState, useRef } from 'react';
import { useParams, Link } from 'react-router-dom';
import { api } from '../lib/api';
import { Build } from '../types';
import { Terminal } from 'lucide-react';

export function BuildDetails() {
  const { id } = useParams<{ id: string }>();
  const [build, setBuild] = useState<Build | null>(null);
  const [logs, setLogs] = useState<string[]>([]);
  const logsEndRef = useRef<HTMLDivElement>(null);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!id) return;

    api.get<Build>(`/builds/${id}`)
      .then(res => setBuild(res.data))
      .catch(console.error);

    // Setup WebSocket
    // Note: Hardcoding WS URL for now, should be dynamic
    const ws = new WebSocket(`ws://localhost:8080/ws/logs/${id}`);
    
    ws.onmessage = (event) => {
      setLogs(prev => [...prev, event.data]);
    };

    wsRef.current = ws;

    return () => {
      ws.close();
    };
  }, [id]);

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);

  if (!build) return <div className="p-8 text-white">Loading build...</div>;

  return (
    <div className="p-8 max-w-7xl mx-auto h-screen flex flex-col">
      <div className="mb-6 flex-shrink-0">
        <div className="flex items-center space-x-2 text-sm text-gray-400 mb-2">
          <Link to="/" className="hover:text-white">Projects</Link>
          <span>/</span>
          <Link to={`/projects/${build.project_id}`} className="hover:text-white">Back to Project</Link>
        </div>
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-white">Build #{build.id.substring(0, 8)}</h1>
          <span className={`px-3 py-1 rounded-full text-sm font-medium
            ${build.status === 'SUCCESS' ? 'bg-green-900 text-green-200' :
              build.status === 'FAILED' ? 'bg-red-900 text-red-200' :
              build.status === 'RUNNING' ? 'bg-blue-900 text-blue-200' :
              'bg-gray-700 text-gray-200'}`}>
            {build.status}
          </span>
        </div>
      </div>

      <div className="flex-1 bg-gray-900 rounded-lg border border-gray-700 font-mono text-sm overflow-hidden flex flex-col shadow-2xl">
        <div className="px-4 py-2 bg-gray-800 border-b border-gray-700 flex items-center text-gray-400">
          <Terminal className="w-4 h-4 mr-2" />
          <span>Build Logs</span>
        </div>
        <div className="flex-1 overflow-auto p-4 space-y-1">
          {logs.map((log, i) => (
            <div key={i} className="text-gray-300 whitespace-pre-wrap">{log}</div>
          ))}
          <div ref={logsEndRef} />
        </div>
      </div>
    </div>
  );
}
