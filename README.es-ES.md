

# Solana Insider Monitor

Una herramienta para monitorear actividades de carteras Solana, detectar cambios en el saldo y recibir alertas en tiempo real.

## Comunidad

Únete a nuestra comunidad de Discord para:
- Obtener ayuda con la instalación y configuración
- Compartir comentarios y sugerencias
- Conectar con otros usuarios
- Mantente actualizado con las nuevas funciones y lanzamientos
- Discutir sobre el desarrollo de Solana

👉 [Únete al Servidor de Discord](https://discord.gg/7vY9ZBPdya)

## Características

- 🔍 Monitoreo simultáneo de múltiples carteras Solana
- 💰 Seguimiento de cambios en el saldo de tokens
- ⚡ Alertas en tiempo real para cambios significativos
- 🔔 Integración con Discord para notificaciones
- 💾 Almacenamiento persistente de datos de carteras
- 🛡️ Manejo adecuado de interrupciones de red

---

## ⚠️ Importante: Configuración del Endpoint RPC

**El problema de configuración más común es usar el endpoint RPC público por defecto**, que tiene límites de tasa estrictos y causará fallos en el escaneo. Siga esta guía para obtener un endpoint RPC adecuado:

### 🚀 Proveedores RPC Recomendados (Planes Gratuitos Disponibles)

| Proveedor | Plan Gratuito | Velocidad | Configuración |
|----------|-----------|-------|-------|
| **Helius** | 100k solicitudes/día | ⚡⚡⚡ | [Obtener Cuenta Gratis](https://helius.dev) |
| **QuickNode** | 30M solicitudes/mes | ⚡⚡⚡ | [Obtener Cuenta Gratis](https://quicknode.com) |
| **Triton** | 10M solicitudes/mes | ⚡⚡ | [Obtener Cuenta Gratis](https://triton.one) |
| **GenesysGo** | Límites personalizados | ⚡⚡ | [Obtener Cuenta](https://genesysgo.com) |

### ❌ Evite Estos (Limitados por Tasa)
```
❌ https://api.mainnet-beta.solana.com (predeterminado - limitado por tasa)
❌ https://api.devnet.solana.com (solo para desarrollo)
❌ https://solana-api.projectserum.com (limitado por tasa)
```

### ✅ Cómo Configurar tu RPC

1. **Regístrate** en cualquiera de los proveedores anteriores (¡son gratuitos!)
2. **Obtén tu URL RPC** desde el panel de control
3. **Actualiza tu config.json**:
   ```json
   {
     "network_url": "https://your-custom-rpc-endpoint.com",
     ...
   }
   ```

---

## Inicio Rápido

### Requisitos Previos

- Go 1.23.2 o posterior
- **Un endpoint RPC dedicado de Solana** (vea [Configuración RPC](#️-important-rpc-endpoint-setup) arriba - ¡esto es crucial!)

### Instalación

```bash
# Clona el repositorio
git clone https://github.com/accursedgalaxy/insider-monitor
cd insider-monitor

# Instala las dependencias
go mod download
```

### Configuración

1. Copia la configuración de ejemplo:
```bash
cp config.example.json config.json
```

2. **⚠️ IMPORTANTE**: Edita `config.json` y reemplaza el endpoint RPC:
```json
{
    "network_url": "YOUR_DEDICATED_RPC_URL_HERE",
    "wallets": [
        "YOUR_WALLET_ADDRESS_1",
        "YOUR_WALLET_ADDRESS_2"
    ],
    "scan_interval": "1m",
    "alerts": {
        "minimum_balance": 1000,
        "significant_change": 0.20,
        "ignore_tokens": []
    },
    "discord": {
        "enabled": false,
        "webhook_url": "",
        "channel_id": ""
    }
}
```

3. **Obtén tu endpoint RPC** de los [proveedores listados arriba](#-recommended-rpc-providers-free-tiers-available) y actualiza `network_url`

### Opciones de Configuración

- `network_url`: **URL de tu endpoint RPC dedicado** (ver la sección de configuración RPC arriba)
- `wallets`: Array de direcciones de carteras Solana para monitorear
- `scan_interval`: Tiempo entre escaneos (ej., "30s", "1m", "5m")
- `alerts`:
  - `minimum_balance`: Saldo mínimo de tokens para activar alertas
  - `significant_change`: Porcentaje de cambio para activar alertas (0.20 = 20%)
  - `ignore_tokens`: Array de direcciones de tokens a ignorar
- `discord`:
  - `enabled`: Establezca en `true` para habilitar notificaciones de Discord
  - `webhook_url`: URL del webhook de Discord
  - `channel_id`: ID del canal de Discord
- `scan`:
  - `scan_mode`: Modo de escaneo de tokens
    - `"all"`: Monitorear todos los tokens (predeterminado)
    - `"whitelist"`: Monitorear solo tokens en `include_tokens`
    - `"blacklist"`: Monitorear todos los tokens excepto los de `exclude_tokens`
  - `include_tokens`: Array de direcciones de tokens para monitorear específicamente (usado con el modo `whitelist`)
  - `exclude_tokens`: Array de direcciones de tokens a ignorar (usado con el modo `blacklist`)

### Ejemplos de Modo de Escaneo

Aquí hay ejemplos de diferentes configuraciones de escaneo:

1. Monitorear todos los tokens:
```json
{
    "scan": {
        "scan_mode": "all",
        "include_tokens": [],
        "exclude_tokens": []
    }
}
```

2. Monitorear solo tokens específicos (lista blanca):
```json
{
    "scan": {
        "scan_mode": "whitelist",
        "include_tokens": [
            "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",  // USDC
            "So11111111111111111111111111111111111111112"     // SOL
        ],
        "exclude_tokens": []
    }
}
```

3. Monitorear todos los tokens excepto unos específicos (lista negra):
```json
{
    "scan": {
        "scan_mode": "blacklist",
        "include_tokens": [],
        "exclude_tokens": [
            "TokenAddressToIgnore1",
            "TokenAddressToIgnore2"
        ]
    }
}
```

### Ejecución del Monitor

```bash
go run cmd/monitor/main.go
```

#### Archivo de Configuración Personalizado
```bash
go run cmd/monitor/main.go -config path/to/config.json
```

### Niveles de Alerta

El monitor utiliza tres niveles de alerta basados en el `significant_change` configurado:
- 🔴 **Crítico**: Cambios >= 5x el umbral
- 🟡 **Advertencia**: Cambios >= 2x el umbral
- 🟢 **Información**: Cambios por debajo de 2x el umbral

### Almacenamiento de Datos

El monitor almacena los datos de las carteras en el directorio `./data` para:
- Prevenir alertas falsas después de los reinicios
- Rastrear cambios históricos
- Manejar adecuadamente las interrupciones de red

### Compilación desde el Código Fuente

```bash
make build
```

El binario estará disponible en el directorio `bin`.

## 🔧 Solución de Problemas

### Problemas Comunes y Soluciones

#### ❌ "Rate limit exceeded" / "Too Many Requests" Error
**Problema**: Uso del endpoint RPC público predeterminado que tiene límites de tasa estrictos
```
❌ Rate limit exceeded after 5 retries
```

**Solución**: 
1. Obtén un endpoint RPC gratuito de [uno de los proveedores anteriores](#-recommended-rpc-providers-free-tiers-available)
2. Actualiza tu `config.json` con el nuevo endpoint:
   ```json
   {
     "network_url": "https://your-custom-rpc-endpoint.com",
     ...
   }
   ```

#### ❌ "Invalid wallet address format" Error
**Problema**: Formato de dirección de cartera incorrecto en config.json
```
❌ invalid wallet address format at index 0: abc123
```

**Solución**: Asegúrese de que las direcciones de la cartera sean direcciones Solana codificadas en base58 válidas (32-44 caracteres)
```json
{
  "wallets": [
    "CvQk2xkXtiMj2JqqVx1YZkeSqQ7jyQkNqqjeNE1jPTfc"  ✅ Valid format
  ]
}
```

#### ❌ "Configuration file not found" Error
**Problema**: config.json no existe
```
❌ Configuration file not found: config.json
```

**Solución**: 
```bash
cp config.example.json config.json
```

#### ❌ "Connection check failed" Error
**Problema**: Problemas de red o del endpoint RPC

**Solución**:
1. Verifique su conexión a internet
2. Confirme que la URL del endpoint RPC sea correcta
3. Intente con un proveedor RPC diferente
4. Pruebe su endpoint RPC manualmente:
   ```bash
   curl -X POST -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","id":1,"method":"getSlot"}' \
     YOUR_RPC_ENDPOINT_URL
   ```

### Obtener Ayuda

Si aún tiene problemas:
1. Consulte nuestra [comunidad de Discord](https://discord.gg/7vY9ZBPdya) para obtener ayuda
2. Revise los registros para mensajes de error específicos
3. Asegúrese de tener la última versión del monitor
4. Intente los pasos de solución de problemas anteriores

## Contribuir

¡Las contribuciones son bienvenidas! No dude en enviar un Pull Request.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - consulte el archivo LICENSE para más detalles.
