import React from "react";
import { Play, Upload, FileText, AlertCircle, Clock, CheckCircle, AlertTriangle, Trash2, TrendingUp } from "lucide-react";

const getStatusColor = (status) => {
  switch (status) {
    case "success":
      return { bg: "bg-green-50", border: "border-green-200", text: "text-green-700", icon: "text-green-600" };
    case "failed":
      return { bg: "bg-red-50", border: "border-red-200", text: "text-red-700", icon: "text-red-600" };
    case "running":
      return { bg: "bg-blue-50", border: "border-blue-200", text: "text-blue-700", icon: "text-blue-600" };
    default:
      return { bg: "bg-gray-50", border: "border-gray-200", text: "text-gray-700", icon: "text-gray-600" };
  }
};

const StatusIcon = ({ status }) => {
  switch (status) {
    case "success":
      return <CheckCircle className="w-5 h-5 text-green-600" />;
    case "failed":
      return <AlertTriangle className="w-5 h-5 text-red-600" />;
    case "running":
      return <div className="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />;
    default:
      return <Clock className="w-5 h-5 text-gray-600" />;
  }
};

export default function WorkflowList({
  workflows,
  runs,
  workflowStats,
  onUpload,
  onTrigger,
  onSelectRun,
  onDeleteWorkflow,
  onViewHistory,
  uploadStatus,
  loading,
}) {
  const getLatestRun = (workflowName) => {
    return runs
      .filter((r) => r.workflow_name === workflowName)
      .sort((a, b) => new Date(b.started_at) - new Date(a.started_at))[0];
  };

  const formatTime = (dateString) => {
    if (!dateString) return "—";
    const date = new Date(dateString);
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  };

  const formatDate = (dateString) => {
    if (!dateString) return "—";
    const date = new Date(dateString);
    return date.toLocaleDateString([], { month: "short", day: "numeric" });
  };

  const formatDuration = (seconds) => {
    if (!seconds) return "—";
    if (seconds < 60) return `${Math.round(seconds)}s`;
    if (seconds < 3600) return `${Math.round(seconds / 60)}m`;
    return `${Math.round(seconds / 3600)}h`;
  };

  return (
    <div className="space-y-6">
      {/* Upload Section */}
      <div className="bg-white rounded-lg border border-gray-200 p-6">
        <div className="flex items-start gap-4">
          <AlertCircle className="w-6 h-6 text-blue-600 flex-shrink-0 mt-1" />
          <div className="flex-1">
            <h3 className="font-semibold text-gray-900 mb-1">
              Add a new workflow
            </h3>
            <p className="text-sm text-gray-600 mb-4">
              Upload a YAML workflow file to automate your CI/CD pipeline.
            </p>
            <div className="flex items-center gap-3">
              <label className="inline-flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-md cursor-pointer transition">
                <Upload className="w-4 h-4 mr-2" />
                Upload workflow file
                <input
                  type="file"
                  accept=".yml,.yaml"
                  onChange={onUpload}
                  className="hidden"
                />
              </label>
              {uploadStatus && (
                <p
                  className={`text-sm font-medium ${
                    uploadStatus.startsWith("✓")
                      ? "text-green-600"
                      : "text-red-600"
                  }`}
                >
                  {uploadStatus}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Workflows Grid */}
      <div>
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Workflows</h2>
        {workflows.length === 0 ? (
          <div className="bg-white rounded-lg border border-gray-200 px-8 py-12 text-center">
            <FileText className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <p className="text-gray-600 text-base mb-2 font-medium">
              No workflows yet
            </p>
            <p className="text-gray-500 text-sm">
              Upload your first workflow file to get started.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {workflows.map((wf) => {
              const latestRun = getLatestRun(wf.name);
              const stats = workflowStats[wf.name] || {};
              const colors = getStatusColor(latestRun?.status);

              return (
                <div
                  key={wf.name}
                  className="bg-white rounded-lg border border-gray-200 overflow-hidden hover:border-gray-300 hover:shadow-md transition"
                >
                  {/* Status Bar */}
                  {latestRun && (
                    <div className={`${colors.bg} ${colors.border} border-b px-4 py-2 flex items-center justify-between`}>
                      <div className="flex items-center gap-2">
                        <StatusIcon status={latestRun.status} />
                        <span className={`text-sm font-medium ${colors.text} capitalize`}>
                          {latestRun.status}
                        </span>
                      </div>
                      <span className="text-xs text-gray-500">
                        {formatTime(latestRun.started_at)}
                      </span>
                    </div>
                  )}

                  {/* Workflow Card Content */}
                  <div className="p-4">
                    {/* Workflow Name & Jobs */}
                    <div className="mb-4">
                      <h3 className="text-base font-semibold text-gray-900 truncate">
                        {wf.name}
                      </h3>
                      <p className="text-xs text-gray-500 mt-1">
                        {Object.keys(wf.jobs || {}).length} job
                        {Object.keys(wf.jobs || {}).length !== 1 ? "s" : ""}
                      </p>
                    </div>

                    {/* Stats Bar */}
                    {stats.total_runs > 0 && (
                      <div className="mb-4 p-3 bg-gray-50 rounded border border-gray-200">
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-2 text-xs">
                            <TrendingUp className="w-4 h-4 text-gray-600" />
                            <span className="text-gray-600">
                              <span className="font-semibold text-gray-900">{stats.total_runs}</span> total runs
                            </span>
                          </div>
                          {stats.success_rate !== undefined && (
                            <span className={`text-xs font-semibold ${stats.success_rate >= 80 ? 'text-green-600' : stats.success_rate >= 50 ? 'text-yellow-600' : 'text-red-600'}`}>
                              {Math.round(stats.success_rate)}%
                            </span>
                          )}
                        </div>
                        <div className="flex items-center gap-2 text-xs text-gray-600">
                          <span>
                            ✓ <span className="font-semibold text-green-600">{stats.successful_runs}</span>
                          </span>
                          {stats.failed_runs > 0 && (
                            <span>
                              ✗ <span className="font-semibold text-red-600">{stats.failed_runs}</span>
                            </span>
                          )}
                          {stats.average_duration > 0 && (
                            <span>
                              ⏱ <span className="font-semibold text-gray-700">{formatDuration(stats.average_duration)}</span>
                            </span>
                          )}
                        </div>
                      </div>
                    )}

                    {/* Latest Run Info */}
                    {latestRun ? (
                      <div className="mb-4 p-3 bg-gray-50 rounded border border-gray-200">
                        <p className="text-xs text-gray-600 mb-2">
                          <span className="font-medium">Latest run:</span> #{latestRun.id.slice(-6)}
                        </p>
                        <p className="text-xs text-gray-600">
                          <span className="font-medium">Triggered:</span> {formatDate(latestRun.started_at)} at {formatTime(latestRun.started_at)}
                        </p>
                      </div>
                    ) : (
                      <div className="mb-4 p-3 bg-gray-50 rounded border border-gray-200">
                        <p className="text-xs text-gray-600">
                          <span className="font-medium">Status:</span> Never run
                        </p>
                      </div>
                    )}

                    {/* Actions */}
                    <div className="flex gap-2">
                      <button
                        onClick={() => onTrigger(wf.name)}
                        disabled={loading}
                        className="flex-1 inline-flex items-center justify-center px-3 py-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white text-sm font-medium rounded transition"
                      >
                        <Play className="w-4 h-4 mr-1" />
                        Run
                      </button>
                      {latestRun && (
                        <button
                          onClick={() => onSelectRun(latestRun.id)}
                          className="flex-1 inline-flex items-center justify-center px-3 py-2 bg-gray-100 hover:bg-gray-200 text-gray-900 text-sm font-medium rounded transition"
                        >
                          Logs
                        </button>
                      )}
                      <button
                        onClick={() => onViewHistory(wf.name)}
                        className="flex-1 inline-flex items-center justify-center px-3 py-2 bg-purple-50 hover:bg-purple-100 text-purple-600 text-sm font-medium rounded transition"
                        title="View all runs"
                      >
                        <Clock className="w-4 h-4 mr-1" />
                        History
                      </button>
                      <button
                        onClick={() => onDeleteWorkflow(wf.name)}
                        className="inline-flex items-center justify-center px-3 py-2 bg-red-50 hover:bg-red-100 text-red-600 text-sm font-medium rounded transition"
                        title="Delete workflow"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}