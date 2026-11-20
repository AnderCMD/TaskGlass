# Code signing

Unsigned programs that hook into other processes are often flagged by antivirus
heuristics and by Windows SmartScreen. Signing `taskglass.exe`,
`taskglass_tap.dll` and the installer with a trusted certificate removes most
of those false positives.

Signing is free for open-source projects through the
[SignPath Foundation](https://signpath.org/) (certificate issued to the
foundation, signing done by [SignPath.io](https://signpath.io/)). The release
workflow is already prepared: it **stays unsigned until you complete the setup
below**, then signs by itself. Nothing else has to change.

[Versión en español más abajo.](#firma-de-código)

## One-time setup

1. **Apply** for the open-source programme on the SignPath Foundation site
   (project: TaskGlass, repository `AnderCMD/TaskGlass`, Apache-2.0). They review
   the project by hand; approval is not instant.
2. In SignPath create the **project** (slug `TaskGlass`), the **signing policy**
   `release-signing`, and two **artifact configurations** (below). Connect the
   project to the GitHub repository (trusted build system: GitHub Actions).
3. Create an **API token** for a CI user that may submit to that policy.
4. In the GitHub repository (**Settings → Secrets and variables → Actions**) add:

   | Kind | Name | Value |
   |------|------|-------|
   | Variable | `SIGNPATH_ORGANIZATION_ID` | your SignPath organization id (turns signing on) |
   | Secret | `SIGNPATH_API_TOKEN` | the CI user's API token |
   | Variable (optional) | `SIGNPATH_PROJECT_SLUG` | default `TaskGlass` |
   | Variable (optional) | `SIGNPATH_SIGNING_POLICY_SLUG` | default `release-signing` |

5. Add the code signing policy statement SignPath asks for to the README, for
   example:

   > Free code signing provided by [SignPath.io](https://signpath.io), certificate
   > by [SignPath Foundation](https://signpath.org). TaskGlass makes no network connections; see
   > `SECURITY.md`.

6. Push a tag (`vX.Y.Z`). The workflow signs the binaries, builds the zip and the
   installer from the signed binaries, signs the installer, checks every
   signature with `Get-AuthenticodeSignature` (the release fails if one is
   invalid) and publishes.

If the signing policy requires manual approval, the workflow waits (up to an
hour) for you to approve the request in SignPath.

## Artifact configurations

`binaries` (the zip with `taskglass.exe` and `taskglass_tap.dll`):

```xml
<artifact-configuration xmlns="http://signpath.io/artifact-configuration/v1">
  <zip-file>
    <pe-file-set>
      <include path="*.exe" />
      <include path="*.dll" />
      <authenticode-sign />
    </pe-file-set>
  </zip-file>
</artifact-configuration>
```

`installer` (the zip with `TaskGlass-<version>-Setup.exe`):

```xml
<artifact-configuration xmlns="http://signpath.io/artifact-configuration/v1">
  <zip-file>
    <pe-file-set>
      <include path="*-Setup.exe" />
      <authenticode-sign />
    </pe-file-set>
  </zip-file>
</artifact-configuration>
```

SignPath validates a configuration when you save it; if its editor suggests a
different element name for your plan, follow the editor.

## Good to know

- Signing does not make an antivirus trust a program instantly: SmartScreen and
  some engines also look at download reputation, which builds up over time.
  A signed file with a stable publisher identity is, however, far less likely to
  be quarantined.
- Never commit certificates or tokens. The workflow only reads the secret above.
- Until signing is enabled, report a false positive to the antivirus vendor and
  verify downloads with `SHA256SUMS.txt`.

---

# Firma de código

Los programas sin firmar que se enganchan a otros procesos suelen ser marcados
por la heurística de los antivirus y por Windows SmartScreen. Firmar
`taskglass.exe`, `taskglass_tap.dll` y el instalador con un certificado de
confianza elimina la mayoría de esos falsos positivos.

La firma es gratuita para proyectos de código abierto mediante la
[SignPath Foundation](https://signpath.org/) (el certificado pertenece a la
fundación y la firma la hace [SignPath.io](https://signpath.io/)). El workflow de
publicación ya está preparado: **se mantiene sin firmar hasta que completes la
configuración**, y entonces firma solo. No hay que cambiar nada más.

## Configuración (una sola vez)

1. **Solicita** el programa de código abierto en el sitio de SignPath Foundation
   (proyecto TaskGlass, repositorio `AnderCMD/TaskGlass`, Apache-2.0). Revisan el
   proyecto a mano; la aprobación no es inmediata.
2. En SignPath crea el **proyecto** (slug `TaskGlass`), la **política de firma**
   `release-signing` y las dos **configuraciones de artefacto** (arriba), y
   conecta el proyecto con el repositorio de GitHub (sistema de compilación de
   confianza: GitHub Actions).
3. Crea un **token de API** para un usuario de CI que pueda enviar a esa política.
4. En el repositorio de GitHub (**Settings → Secrets and variables → Actions**)
   añade la variable `SIGNPATH_ORGANIZATION_ID` (activa la firma) y el secreto
   `SIGNPATH_API_TOKEN`. Opcionales: las variables `SIGNPATH_PROJECT_SLUG` y
   `SIGNPATH_SIGNING_POLICY_SLUG` si usas otros nombres.
5. Añade al README la declaración de política de firma que pide SignPath (el
   ejemplo de arriba, traducida si quieres).
6. Sube un tag (`vX.Y.Z`). El workflow firma los binarios, construye el zip y el
   instalador con los binarios firmados, firma el instalador, comprueba cada
   firma con `Get-AuthenticodeSignature` (la release falla si alguna no es
   válida) y publica.

Si la política exige aprobación manual, el workflow espera (hasta una hora) a que
apruebes la solicitud en SignPath.

## Conviene saber

- Firmar no hace que un antivirus confíe al instante: SmartScreen y algunos
  motores también valoran la reputación de la descarga, que se gana con el
  tiempo. Aun así, un archivo firmado con un editor estable tiene muchas menos
  probabilidades de ir a cuarentena.
- Nunca subas certificados ni tokens al repositorio. El workflow solo lee el
  secreto indicado.
- Mientras la firma no esté activa, notifica el falso positivo al fabricante del
  antivirus y verifica las descargas con `SHA256SUMS.txt`.
