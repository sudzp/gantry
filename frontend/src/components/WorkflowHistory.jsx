import React, { useState, useEffect } from "react";
import {
  ArrowLeft,
  CheckCircle2,
  XCircle,
  Clock,
  Loader2,
} from "lucide-react";
import apiService from "../services/apiService";

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

const formatDate = (dateString) => {
  if (!dateString) return "—";
  const date = new Date(dateString);
  return date.toLocaleDateString([], { 
    month: "short", 
    day: "numeric", 
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
};

const formatDuration = (seconds) => {
  if (!seconds) return "—";
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  if (mins > 0) return `${mins}m ${secs}s`;
  return `${secs}s`;
};

export default function WorkflowHistory({ workflowName, onBack, onSelectRun }) {
  const [runs, setRuns] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRuns = async () => {
      try {
        setLoading(true);
        const data = await apiService.getWorkflowRuns(workflowName);
        setRuns(data || []);
      } catch (err) {
        console.error("Failed to fetch workflow runs:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchRuns();
  }, [workflowName]);

  const calculateDuration = (run) => {
    if (!run.started_at || !run.completed_at) return null;
    const start = new Date(run.started_at);
    const end = new Date(run.completed_at);
    return Math.floor((end - start) / 1000);
  };

  return (
    <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
      {/* Header */}
      <div className="px-6 py-4 border-b border-gray-200 bg-gray-50 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <button
            onClick={onBack}
            className="p-1 text-gray-600 hover:bg-gray-200 rounded transition"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <div>
            <h2 className="text-lg font-semibold text-gray-900">
              Run History: {workflowName}
            </h2>
            <p className="text-sm text-gray-500 mt-1">
              {runs.length} {runs.length === 1 ? "run" : "runs"}
            </p>
          </div>
        </div>
      </div>

      {/* Runs Table */}
      <div className="overflow-x-auto">
        {loading ? (
          <div className="px-6 py-12 text-center">
            <Loader2 className="w-8 h-8 text-gray-400 mx-auto mb-3 animate-spin" />
            <p className="text-gray-500">Loading runs...</p>
          </div>
        ) : runs.length === 0 ? (
          <div className="px-6 py-12 text-center">
            <Clock className="w-12 h-12 text-gray-300 mx-auto mb-3" />
            <p className="text-gray-500">No runs for this workflow yet</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">
                  Run ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">
                  Started
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">
                  Duration
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-700 uppercase">
                  Action
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {runs.map((run) => {
                const duration = calculateDuration(run);
                return (
                  <tr key={run.id} className="hover:bg-gray-50 transition">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(run.status)}
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getStatusColor(
                            run.status
                          )}`}
                        >
                          {run.status.charAt(0).toUpperCase() + run.status.slice(1)}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-700 font-mono">
                      {run.id.substring(0, 12)}...
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-700">
                      {formatDate(run.started_at)}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-700">
                      {duration ? formatDuration(duration) : "—"}
                    </td>
                    <td className="px-6 py-4">
                      <button
                        onClick={() => onSelectRun(run.id)}
                        className="text-sm px-3 py-1 bg-blue-50 text-blue-700 border border-blue-200 rounded hover:bg-blue-100 transition"
                      >
                        View Logs
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
