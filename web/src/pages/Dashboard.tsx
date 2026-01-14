import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { Project } from '../types';
import { Plus, GitBranch, Clock } from 'lucide-react';
import { Link } from 'react-router-dom';
import { format } from 'date-fns';

export function Dashboard() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // In a real app, we'd handle auth state better.
    // Assuming backend will trust us or we have a cookie now.
    api.get<Project[]>('/projects')
      .then(res => setProjects(res.data))
      .catch(err => console.error("Failed to fetch projects", err))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="p-8 text-white">Loading projects...</div>;
  }

  return (
    <div className="p-8 max-w-7xl mx-auto">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-white">Projects</h1>
        <button className="flex items-center px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md transition-colors">
          <Plus className="w-4 h-4 mr-2" />
          New Project
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {projects.map((project) => (
          <Link 
            key={project.id} 
            to={`/projects/${project.id}`}
            className="block bg-gray-800 rounded-lg p-6 border border-gray-700 hover:border-indigo-500 transition-colors"
          >
            <div className="flex items-start justify-between">
              <div>
                <h3 className="text-xl font-semibold text-white mb-2">{project.name}</h3>
                <p className="text-gray-400 text-sm mb-4">{project.repo_url}</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-4 text-sm text-gray-400">
              <div className="flex items-center">
                <GitBranch className="w-4 h-4 mr-1" />
                {project.default_branch}
              </div>
              <div className="flex items-center">
                <Clock className="w-4 h-4 mr-1" />
                {format(new Date(project.updated_at), 'MMM d, yyyy')}
              </div>
            </div>
          </Link>
        ))}
        
        {projects.length === 0 && (
          <div className="col-span-full text-center py-12 text-gray-500">
            No projects found. Create one to get started!
          </div>
        )}
      </div>
    </div>
  );
}
