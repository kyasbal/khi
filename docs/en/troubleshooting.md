# 🛠 Troubleshooting

## KHI Timeline does not render correctly or at all

KHI relies on WebGL for rendering its timeline and diagrams. In some environments, browser restrictions may block WebGL, causing graphics not to display correctly.

### Browser Restrictions (Chrome 139+ without a GPU)

Starting with Chrome 139, CPU-based WebGL emulation is disabled by default in environments without a GPU.
If you are running KHI on a remote server without a GPU and connecting via Remote Desktop, timelines may not render correctly.

To resolve this issue, you can override the setting in Chrome:

1. Open a new tab and navigate to `chrome://flags/#ignore-gpu-blocklist`.
2. Change the setting for this flag to `Enabled`.
3. Relaunch the browser.

> **Note:** This does not affect most consumer laptops (even those without dedicated graphics cards), as they typically rely on integrated GPUs (iGPUS).

---

> KHI is only tested on the latest version of Google Chrome. It may work with other browsers, but we do not provide support if it does not.
