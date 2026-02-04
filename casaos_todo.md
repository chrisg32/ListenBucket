# Publishing to the CasaOS Built-in App Store

To get your app into the **official CasaOS App Store**, you need to submit a Pull Request to the [IceWhaleTech/CasaOS-AppStore](https://github.com/IceWhaleTech/CasaOS-AppStore) repository.

## Steps

1. **Test thoroughly** on your own CasaOS instance first (mandatory)
2. **Fork** the repository
3. **Create your app directory** with the required files
4. **Submit a Pull Request** and assign it to a CasaOS team member

## Required Files

Your app directory must contain:

| File | Required | Description |
|------|----------|-------------|
| `docker-compose.yml` | Yes | Valid Docker Compose config with `x-casaos` metadata |
| `icon.png` | Yes | App icon (192×192 for featured apps) |
| `screenshot-1.png` | Yes | At least one screenshot showing it working |
| `thumbnail.png` | Optional | Only for featured apps (784×442) |

## Docker Compose Requirements

- App `name` must match pattern `^[a-z0-9][a-z0-9_-]*$`
- Use **specific image tags** (e.g., `:1.2.3`), never `:latest`
- Include `x-casaos` metadata with:
  - `architectures` (amd64, arm, arm64)
  - `main` (primary service name)
  - `author`, `developer`, `category`
  - `title`, `description`, `tagline` (supports locales like `en_us`, `zh_cn`)
  - `icon`, `port_map`
- Use system variables: `$PGID`, `$PUID`, `$TZ`, `$AppID`

## Alternative: Create Your Own App Store

If you want to host apps yourself, you can create a **custom third-party app store** (a Git repo with compose files) and add it via the CasaOS dashboard's "Add Source" button.

## Resources

- [CasaOS-AppStore GitHub](https://github.com/IceWhaleTech/CasaOS-AppStore)
- [Create Custom AppStore Guide](https://awesome.casaos.io/content/3rd-party-app-stores/create-your-first-custom-appstore.html)
