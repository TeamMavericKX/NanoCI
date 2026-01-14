import { Github } from 'lucide-react';

export function LoginPage() {
  const handleLogin = () => {
    window.location.href = 'http://localhost:8080/auth/login';
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white">
      <div className="max-w-md w-full space-y-8 p-8 bg-gray-800 rounded-lg shadow-lg text-center">
        <div>
          <h2 className="mt-6 text-3xl font-extrabold">Sign in to NanoCI</h2>
          <p className="mt-2 text-sm text-gray-400">
            Your private, container-native CI/CD platform
          </p>
        </div>
        <button
          onClick={handleLogin}
          className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-gray-900 bg-white hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          <span className="absolute left-0 inset-y-0 flex items-center pl-3">
            <Github className="h-5 w-5 text-gray-900" aria-hidden="true" />
          </span>
          Sign in with GitHub
        </button>
      </div>
    </div>
  );
}
