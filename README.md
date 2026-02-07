# shadowpkg
A tool to download the latest Shadow Tech software packages, created after my Shadow PC experienced some software issues.

> [!NOTE]
> This tool is not affiliated with, authorized, maintained, sponsored, or endorsed by Shadow (Shadow Tech) or any of its affiliates or subsidiaries.
>
> This is an unofficial, community-driven project. Any software packages downloaded via this tool are sourced from publicly available official repositories, but the tool itself is entirely separate from the official Shadow PC development team.
>
> The name "Shadow," "Shadow PC," and related logos are registered trademarks of Shadow Tech. These terms are used here solely for descriptive purposes to identify the service for which this tool is intended.
>
> Use of this tool is at your own risk. The developer of shadowpkg is not responsible for any issues, account actions, or technical failures that may arise from using this software or the packages it retrieves. Always ensure you are complying with Shadow's official Terms of Service.

## Usage
A pre-compiled Windows binary can be found in the [releases tab.](https://github.com/dwilliamsuk/shadowpkg/releases/tag/v0.0.1)

```
Usage of shadowpkg:
  -environment string
        The environment to download packages from (default "prod")
  -output string
        The output directory to download packages to (default ".")
```

### Example Usage
```
PS C:\Users\Shadow\Desktop> .\shadowpkg.exe -environment="prod" -output=".\output_dir"
Downloading [15] packages...
Done! Saved packages to [C:\Users\Shadow\Desktop\output_dir]
```

> [!CAUTION]
> This tool <ins>**WILL OVERWRITE FILES**</ins> with the same name in the output directory. Ensure that the output directory specified does not contain important information.

## Documentation
> [!NOTE]
> This tool was created by reverse engineering tooling available on my own Shadow PC. As such, I can only provide my own notes as documentation on the software package retrieval process.

### Known Environments
- ``prod`` - The regular production software used on a Shadow PC
- ``emergency`` - The recovery software packages pulled when required on a Shadow PC

### Known Endpoints

#### ``https://oblivion.shadow.tech/environments``
Contains the file hashes / IDs of the software package manifests for a particular environment.

For example, the production manifest file ID can be found at:

``https://oblivion.shadow.tech/environments/prod``

-----

#### ``https://oblivion.shadow.tech/manifests``
Contains the software package manifest files.

For example, the current production manifest file can be found at:

``https://oblivion.shadow.tech/manifests/a3c0a60370a923197dc441c38a1c7b55570fd5eb``

-----

#### ``https://oblivion.shadow.tech/packages``
Contains the actual software packages themselves, along with MD5 hashes.

For example, the ``ShadowStreamer`` software can be found at:

``https://oblivion.shadow.tech/packages/ShadowStreamer/6.0.6/ShadowStreamer.msi``

with it's MD5 hash at:

``https://oblivion.shadow.tech/packages/ShadowStreamer/6.0.6/ShadowStreamer.md5``
