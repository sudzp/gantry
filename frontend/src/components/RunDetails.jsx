import React, { useState } from "react";
import {
  CheckCircle2,
  XCircle,
  Loader2,
  Clock,
  ChevronDown,
  ChevronRight,
  Calendar,
  GitBranch,
  User,
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

function JobItem({ name, job }) {
  const [isExpanded, setIsExpanded] = useState(true);

  const duration =
    job.started_at && job.ended_at
      ? Math.round((new Date(job.ended_at) - new Date(job.started_at)) / 1000)
      : null;

  return (
    <div className="border border-gray-200 rounded-lg overflow-hidden">
      {/* Job Header */}
      <div
        className="px-4 py-3 bg-gray-50 cursor-pointer hover:bg-gray-100 transition"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {isExpanded ? (
              <ChevronDown className="w-4 h-4 text-gray-500" />
            ) : (
              <ChevronRight className="w-4 h-4 text-gray-500" />
            )}
            {getStatusIcon(job.status)}
            <span className="font-semibold text-gray-900">{name}</span>
          </div>
          <div className="flex items-center gap-3">
            {duration !== null && (
              <span className="text-sm text-gray-600">{duration}s</span>
            )}
            <span
              className={`px-2.5 py-1 text-xs font-medium rounded-full border ${getStatusColor(job.status)}`}
            >
              {job.status}
            </span>
          </div>
        </div>
      </div>

      {/* Job Details */}
      {isExpanded && (
        <div className="px-4 py-4 space-y-4">
          {/* Timestamps */}
          {(job.started_at || job.ended_at) && (
            <div className="text-sm text-gray-600 space-y-1">
              {job.started_at && (
                <div className="flex items-center gap-2">
                  <Calendar className="w-4 h-4" />
                  <span>
                    Started: {new Date(job.started_at).toLocaleString()}
                  </span>
                </div>
              )}
              {job.ended_at && (
                <div className="flex items-center gap-2">
                  <Calendar className="w-4 h-4" />
                  <span>Ended: {new Date(job.ended_at).toLocaleString()}</span>
                </div>
              )}
            </div>
          )}

          {/* Steps */}
          {job.steps && job.steps.length > 0 && (
            <div>
              <h4 className="text-sm font-semibold text-gray-900 mb-2">
                Steps
              </h4>
              <div className="space-y-2">
                {job.steps.map((step, idx) => (
                  <div
                    key={idx}
                    className="flex items-center gap-2 text-sm text-gray-700 py-1"
                  >
                    <div className="w-1.5 h-1.5 bg-gray-400 rounded-full"></div>
                    {step.name}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Output Logs */}
          {job.output && (
            <div>
              <h4 className="text-sm font-semibold text-gray-900 mb-2">
                Build logs
              </h4>
              <div className="bg-gray-900 rounded-lg overflow-hidden">
                <div className="px-4 py-2 border-b border-gray-700 flex items-center justify-between">
                  <span className="text-xs text-gray-400 font-mono">
                    Console output
                  </span>
                </div>
                <pre className="px-4 py-3 text-sm text-gray-300 font-mono overflow-x-auto max-h-96">
                  {job.output}
                </pre>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default function RunDetails({ run }) {
  if (!run) return null;

  const renderJobs = () => {
    const jobsToRender =
      run.job_order && run.job_order.length > 0
        ? run.job_order
            .map((jobName) => ({ name: jobName, job: run.jobs[jobName] }))
            .filter((j) => j.job)
        : Object.entries(run.jobs || {}).map(([name, job]) => ({ name, job }));

    return jobsToRender.map(({ name, job }) => (
      <JobItem key={name} name={name} job={job} />
    ));
  };

  return (
    <div className="space-y-6">
      {/* Run Header */}
      <div className="bg-white border border-gray-200 rounded-lg p-6">
        <div className="flex items-start gap-3 mb-4">
          {getStatusIcon(run.status)}
          <div className="flex-1">
            <h2 className="text-xl font-semibold text-gray-900 mb-2">
              {run.workflow_name}
            </h2>
            <div className="flex items-center gap-4 text-sm text-gray-600">
              <span className="flex items-center gap-1">
                <GitBranch className="w-4 h-4" />
                main
              </span>
              <span className="flex items-center gap-1">
                <User className="w-4 h-4" />
                system
              </span>
              <span className="flex items-center gap-1">
                <Calendar className="w-4 h-4" />
                {new Date(run.started_at).toLocaleString()}
              </span>
            </div>
          </div>
          <span
            className={`px-3 py-1.5 text-sm font-medium rounded-full border ${getStatusColor(run.status)}`}
          >
            {run.status}
          </span>
        </div>
        <div className="text-sm text-gray-600">
          <span className="font-mono text-xs bg-gray-100 px-2 py-1 rounded">
            #{run.id.slice(0, 7)}
          </span>
        </div>
      </div>

      {/* Jobs Section */}
      <div className="bg-white border border-gray-200 rounded-lg">
        <div className="px-6 py-4 border-b border-gray-200 bg-gray-50">
          <h3 className="text-lg font-semibold text-gray-900">Jobs</h3>
        </div>
        <div className="p-6 space-y-4">{renderJobs()}</div>
      </div>
    </div>
  );
}