import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { api } from '../lib/api';
import { Project, Build } from '../types';
import { format } from 'date-fns';
import { CheckCircle, XCircle, Clock, Loader, PlayCircle } from 'lucide-react';
import { cn } from '../lib/utils';

export function ProjectDetails() {
  const { id } = useParams<{ id: string }>();
  const [project, setProject] = useState<Project | null>(null);
  const [builds, setBuilds] = useState<Build[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;

    Promise.all([
      api.get<Project>(`/projects/${id}`),
      api.get<Build[]>(`/projects/${id}/builds`)
    ])
    .then(([pRes, bRes]) => {
      setProject(pRes.data);
      setBuilds(bRes.data);
    })
    .catch(console.error)
    .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="p-8 text-white">Loading...</div>;
  if (!project) return <div className="p-8 text-white">Project not found</div>;

  return (
    <div className="p-8 max-w-7xl mx-auto">
      <div className="mb-8">
        <div className="flex items-center space-x-2 text-sm text-gray-400 mb-2">
          <Link to="/" className="hover:text-white">Projects</Link>
          <span>/</span>
          <span className="text-white">{project.name}</span>
        </div>
        <h1 className="text-3xl font-bold text-white">{project.name}</h1>
      </div>

      <div className="bg-gray-800 rounded-lg overflow-hidden border border-gray-700">
        <div className="px-6 py-4 border-b border-gray-700 bg-gray-900/50">
          <h2 className="text-lg font-medium text-white">Build History</h2>
        </div>
        <div className="divide-y divide-gray-700">
          {builds.map((build) => (
            <Link 
              key={build.id}
              to={`/builds/${build.id}`}
              className="block p-4 hover:bg-gray-700/50 transition-colors"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <StatusIcon status={build.status} />
                  <div>
                    <div className="text-white font-medium">{build.commit_message || 'No commit message'}</div>
                    <div className="flex items-center space-x-3 text-sm text-gray-400 mt-1">
                      <span className="font-mono text-xs bg-gray-700 px-2 py-0.5 rounded text-gray-300">
                        {build.commit_hash.substring(0, 7)}
                      </span>
                      <span>{build.branch}</span>
                    </div>
                  </div>
                </div>
                <div className="text-sm text-gray-400">
                  {build.finished_at 
                    ? format(new Date(build.finished_at), 'MMM d, HH:mm') 
                    : 'In Progress'}
                </div>
              </div>
            </Link>
          ))}
          {builds.length === 0 && (
            <div className="p-8 text-center text-gray-500">No builds yet.</div>
          )}
        </div>
      </div>
    </div>
  );
}

function StatusIcon({ status }: { status: string }) {
  switch (status) {
    case 'SUCCESS':
      return <CheckCircle className="w-5 h-5 text-green-500" />;
    case 'FAILED':
      return <XCircle className="w-5 h-5 text-red-500" />;
    case 'RUNNING':
      return <Loader className="w-5 h-5 text-blue-500 animate-spin" />;
    case 'PENDING':
      return <Clock className="w-5 h-5 text-yellow-500" />;
    default:
      return <PlayCircle className="w-5 h-5 text-gray-500" />;
  }
}
