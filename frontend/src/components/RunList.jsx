import React from "react";
import {
  RefreshCw,
  CheckCircle2,
  XCircle,
  Clock,
  Loader2,
  GitBranch,
  Calendar,
} from "lucide-react";

const getStatusIcon = (status) => {
  switch (status) {
    case "success":
      return <CheckCircle2 className="w-4 h-4 text-green-600" />;
    case "failed":
      return <XCircle className="w-4 h-4 text-red-600" />;
    case "running":
      return <Loader2 className="w-4 h-4 text-yellow-600 animate-spin" />;
    default:
      return <Clock className="w-4 h-4 text-gray-400" />;
  }
};

const getStatusColor = (status) => {
  switch (status) {
    case "success":
      return "bg-green-50 text-green-700 border-green-200";
    case "failed":
      return "bg-red-50 text-red-700 border-red-200";
    case "running":
      return "bg-yellow-50 text-yellow-700 border-yellow-200";
    default:
      return "bg-gray-50 text-gray-700 border-gray-200";
  }
};

const formatRelativeTime = (date) => {
  const now = new Date();
  const then = new Date(date);
  const diff = now - then;
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (seconds < 60) return "just now";
  if (minutes < 60) return `${minutes}m ago`;
  if (hours < 24) return `${hours}h ago`;
  return `${days}d ago`;
};

export default function RunList({ runs, onSelectRun, onRefresh }) {
  return (
    <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
      {/* Header */}
      <div className="px-6 py-4 border-b border-gray-200 bg-gray-50 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-gray-900">Workflow runs</h2>
        <button
          onClick={onRefresh}
          className="p-2 text-gray-600 hover:bg-gray-100 rounded-md transition"
          title="Refresh runs"
        >
          <RefreshCw className="w-4 h-4" />
        </button>
      </div>

      {/* Runs List */}
      <div className="divide-y divide-gray-200">
        {runs.length === 0 ? (
          <div className="px-6 py-12 text-center">
            <Clock className="w-12 h-12 text-gray-300 mx-auto mb-3" />
            <p className="text-gray-500 text-sm">
              No workflow runs yet. Trigger a workflow to see results here.
            </p>
          </div>
        ) : (
          runs.map((run) => (
            <div
              key={run.id}
              onClick={() => onSelectRun(run.id)}
              className="px-6 py-4 hover:bg-gray-50 cursor-pointer transition"
            >
              <div className="flex items-start gap-3">
                {getStatusIcon(run.status)}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <span className="font-medium text-gray-900 truncate">
                      {run.workflow_name}
                    </span>
                    <span
                      className={`px-2 py-0.5 text-xs font-medium rounded-full border ${getStatusColor(run.status)}`}
                    >
                      {run.status}
                    </span>
                  </div>
                  <div className="flex items-center gap-4 text-sm text-gray-600">
                    <span className="flex items-center gap-1">
                      <GitBranch className="w-3.5 h-3.5" />
                      main
                    </span>
                    <span className="flex items-center gap-1">
                      <Calendar className="w-3.5 h-3.5" />
                      {formatRelativeTime(run.started_at)}
                    </span>
                    <span className="text-gray-400">•</span>
                    <span className="text-xs text-gray-500 font-mono">
                      #{run.id.slice(0, 7)}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}