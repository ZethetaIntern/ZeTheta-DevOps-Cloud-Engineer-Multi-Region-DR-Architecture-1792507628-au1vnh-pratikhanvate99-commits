param(
  [int]$Port = $(if ($env:PAYSECURE_HOST_PORT) { [int]$env:PAYSECURE_HOST_PORT } else { 8080 })
)
$ErrorActionPreference = 'Stop'
$base = "http://localhost:$Port"

function Get-StatusCode([scriptblock]$Action) {
  try { & $Action | Out-Null; return 200 }
  catch {
    if ($_.Exception.Response -and $_.Exception.Response.StatusCode) {
      return [int]$_.Exception.Response.StatusCode
    }
    throw
  }
}

$ready = $false
for ($attempt = 0; $attempt -lt 30; $attempt++) {
  try {
    if ((Invoke-WebRequest -Uri "$base/health/ready" -TimeoutSec 3).StatusCode -eq 200) {
      $ready = $true
      break
    }
  } catch { Start-Sleep -Seconds 1 }
}
if (-not $ready) { throw "API did not become ready at $base" }

$merchantHeaders = @{
  'X-Merchant-ID' = 'merchant-demo-001'
  'X-API-Key' = 'local-demo-api-key'
}
$idempotencyKey = [guid]::NewGuid().ToString()
$headers = $merchantHeaders.Clone()
$headers['Idempotency-Key'] = $idempotencyKey
$body = '{"amount_minor":159900,"currency":"INR"}'
$first = Invoke-RestMethod -Uri "$base/v1/payments" -Method Post -Headers $headers -ContentType 'application/json' -Body $body
$second = Invoke-RestMethod -Uri "$base/v1/payments" -Method Post -Headers $headers -ContentType 'application/json' -Body $body
if ($first.payment_id -ne $second.payment_id) { throw 'Idempotent replay returned a different payment ID' }

$payment = $null
for ($attempt = 0; $attempt -lt 30; $attempt++) {
  Start-Sleep -Milliseconds 500
  $payment = Invoke-RestMethod -Uri "$base/v1/payments/$($first.payment_id)" -Headers $merchantHeaders
  if ($payment.state -eq 'SIMULATED_APPROVED') { break }
}
if ($payment.state -ne 'SIMULATED_APPROVED') { throw "Expected SIMULATED_APPROVED, got $($payment.state)" }

$conflictHeaders = $headers.Clone()
$conflictBody = '{"amount_minor":159901,"currency":"INR"}'
$conflict = Get-StatusCode { Invoke-RestMethod -Uri "$base/v1/payments" -Method Post -Headers $conflictHeaders -ContentType 'application/json' -Body $conflictBody }
if ($conflict -ne 409) { throw "Expected 409 for idempotency conflict, got $conflict" }

$unauthorized = Get-StatusCode { Invoke-RestMethod -Uri "$base/v1/payments" -Method Post -ContentType 'application/json' -Body $body }
if ($unauthorized -ne 401) { throw "Expected 401 for missing credentials, got $unauthorized" }

$invalidHeaders = $merchantHeaders.Clone()
$invalidHeaders['Idempotency-Key'] = [guid]::NewGuid().ToString()
$invalid = Get-StatusCode { Invoke-RestMethod -Uri "$base/v1/payments" -Method Post -Headers $invalidHeaders -ContentType 'application/json' -Body '{"amount_minor":0,"currency":"INR"}' }
if ($invalid -ne 400) { throw "Expected 400 for invalid amount, got $invalid" }

Write-Output "PASS ready=200 payment=$($first.payment_id) state=$($payment.state) replay=same-id conflict=$conflict unauthorized=$unauthorized invalid-amount=$invalid"
