# ⚙️ Settings for Managed Environments

## Google Cloud

### Permissions

The following permissions are required or recommended.

- **Required**
  - `logging.logEntries.list`
- **Recommended**
  - These permissions are used to fetch autocomplete candidates in the New Inspection dialog. KHI works without these permissions, but cluster name suggestions will not be displayed.
    - `monitoring.timeSeries.list`
    - `container.clusters.list` (Only when using Cloud Composer features)

#### Setup

- **Running KHI on environments with a service account attached** (e.g., Google Cloud Compute Engine Instance): Apply the permissions above to the attached service account.
- **Running KHI locally or on Cloud Shell with a user account**: Apply the permissions above to your user account.

### Audit Logging

- **No required configuration** (KHI fully works with the default audit logging configuration).
- **Recommended**
  - Kubernetes Engine API Data access audit logs for `DATA_WRITE`

> **💡 TIP**
> Enabling these will log every patch requests on Pod or Node `.status` field.
> KHI will use this to display detailed container status.
> KHI will still guess the last container status from the audited Pod deletion log even without these logs, however it requires the Pod to be deleted within the queried timeframe.

#### Setup

1. In the Google Cloud Console, [go to the Audit Logs](https://console.cloud.google.com/iam-admin/audit) page.
2. In the Data Access audit logs configuration table, select `Kubernetes Engine API` from the Service column.
3. In the Log Types tab, select the `Data write` Data Access audit log type.
4. Click "SAVE".
