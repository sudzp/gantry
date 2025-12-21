import React from 'react';
import { RefreshCw, CheckCircle, XCircle, Clock, Loader } from 'lucide-react';

const getStatusIcon = (status) => {
  switch (status) {
    case 'success':
      return <CheckCircle className="w-5 h-5 text-green-500" />;
    case 'failed':
      return <XCircle className="w-5 h-5 text-red-500" />;
    case 'running':
      return <Loader className="w-5 h-5 text-blue-500 animate-spin" />;
    default:
      return <Clock className="w-5 h-5 text-gray-400" />;
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case 'success':
      return 'bg-green-100 text-green-800';
    case 'failed':
      return 'bg-red-100 text-red-800';
    case 'running':
      return 'bg-blue-100 text-blue-800';
    default:
      return 'bg-gray-100 text-gray-800';
  }
};

export default function RunList({ runs, onSelectRun, onRefresh }) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold">Recent Runs</h2>
        <button
          onClick={onRefresh}
          className="p-2 text-gray-600 hover:bg-gray-100 rounded transition"
          title="Refresh runs"
        >
          <RefreshCw className="w-5 h-5" />
        </button>
      </div>

      <div className="space-y-3">
        {runs.length === 0 ? (
          <p className="text-gray-500 text-sm">
            No runs yet. Trigger a workflow to see runs here.
          </p>
        ) : (
          runs.map((run) => (
            <div
              key={run.id}
              onClick={() => onSelectRun(run.id)}
              className="p-4 border rounded-lg hover:bg-gray-50 cursor-pointer transition"
            >
              <div className="flex items-center justify-between mb-2">
                <div className="flex items-center gap-3">
                  {getStatusIcon(run.status)}
                  <div>
                    <h3 className="font-semibold">{run.workflow_name}</h3>
                    <p className="text-sm text-gray-500">{run.id}</p>
                  </div>
                </div>
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(run.status)}`}>
                  {run.status}
                </span>
              </div>
              <div className="text-xs text-gray-500">
                Started: {new Date(run.started_at).toLocaleString()}
                {run.completed_at && ` • Completed: ${new Date(run.completed_at).toLocaleString()}`}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}