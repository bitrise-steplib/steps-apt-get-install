# apt-get install

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-apt-get-install?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-apt-get-install/releases)

Install or upgrade packages on Ubuntu with apt-get.

<details>
<summary>Description</summary>

[This Step](https://www.bitrise.io/integrations/steps/apt-get-install) integrates with the Advanced Package Tool (APT) command line tool to install and upgrade packages.

### Configuring the Step
1. Add the **Name of the packages to install/upgrade, separated with spaces**.
2. Add flags to pass on to the `apt-get` command in the **Options for apt-get install/upgrade** input.
3. Allow upgrades to previously installed packages with the **Upgrade packages if previously installed** input.
Under **Options**:
4. Set the level of cache in the **Cache level** input.

### Useful links
- [Installing any additional tools](https://devcenter.bitrise.io/tips-and-tricks/install-additional-tools/#apt-get-on-linux)

### Related Steps
- [Brew Install](https://www.bitrise.io/integrations/steps/brew-install)
- [Flutter Install](https://www.bitrise.io/integrations/steps/flutter-installer)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `packages` | Name of the packages to install/upgrade, separated with spaces | required |  |
| `options` | Flags to pass to the `apt-get install/upgrade` command.  `apt-get install/upgrade -y [options] [packages]`  |  |  |
| `upgrade` | If set to `"yes"`, the step will upgrade the defined packages by calling `apt-get upgrade -y [options] [packages]` command.  Otherwise the step calls `apt-get install -y [options] [packages]`  |  | `yes` |
| `cache_level` | Sets the level of cache.  'all' enables the caching of all files under /var/cache/apt/archives folder (default directory for apt cache). 'none' disables the caching for the step. | required | `all` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-apt-get-install/pulls) and [issues](https://github.com/bitrise-steplib/steps-apt-get-install/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
