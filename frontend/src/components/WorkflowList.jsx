import React from "react";
import { Play, Upload, FileText, AlertCircle } from "lucide-react";

export default function WorkflowList({
  workflows,
  onUpload,
  onTrigger,
  uploadStatus,
  loading,
}) {
  return (
    <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
      {/* Header */}
      <div className="px-6 py-4 border-b border-gray-200 bg-gray-50">
        <h2 className="text-lg font-semibold text-gray-900">All workflows</h2>
      </div>

      {/* Upload Section */}
      <div className="px-6 py-4 border-b border-gray-200 bg-blue-50">
        <div className="flex items-start gap-3">
          <AlertCircle className="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
          <div className="flex-1">
            <h3 className="font-medium text-blue-900 text-sm mb-1">
              Get started with workflows
            </h3>
            <p className="text-sm text-blue-800 mb-3">
              Upload a YAML workflow file to automate your CI/CD pipeline.
            </p>
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
                className={`mt-2 text-sm ${
                  uploadStatus.startsWith("✓")
                    ? "text-green-700"
                    : "text-red-700"
                }`}
              >
                {uploadStatus}
              </p>
            )}
          </div>
        </div>
      </div>

      {/* Workflows List */}
      <div className="divide-y divide-gray-200">
        {workflows.length === 0 ? (
          <div className="px-6 py-12 text-center">
            <FileText className="w-12 h-12 text-gray-300 mx-auto mb-3" />
            <p className="text-gray-500 text-sm">
              No workflows yet. Upload your first workflow to get started.
            </p>
          </div>
        ) : (
          workflows.map((wf) => (
            <div
              key={wf.name}
              className="px-6 py-4 hover:bg-gray-50 transition"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3 flex-1">
                  <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center">
                    <FileText className="w-4 h-4 text-gray-600" />
                  </div>
                  <div className="flex-1">
                    <div className="font-medium text-gray-900">{wf.name}</div>
                    <div className="text-sm text-gray-500 mt-0.5">
                      {Object.keys(wf.jobs || {}).length} job
                      {Object.keys(wf.jobs || {}).length !== 1 ? "s" : ""}
                    </div>
                  </div>
                </div>
                <button
                  onClick={() => onTrigger(wf.name)}
                  disabled={loading}
                  className="inline-flex items-center px-4 py-2 bg-green-600 hover:bg-green-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white text-sm font-medium rounded-md transition"
                  title="Run workflow"
                >
                  <Play className="w-4 h-4 mr-2" />
                  Run workflow
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}