$key = 'ntn_z13943653061c5sQPOZJLETAERl4sUiX4pzTUvF5eGebtA'
$headers = @{
    'Authorization' = "Bearer $key"
    'Notion-Version' = '2022-06-28'
    'Content-Type' = 'application/json'
}
try {
    $r = Invoke-RestMethod -Uri 'https://api.notion.com/v1/pages/3220531c537a80aa8333c42da4ce362f' -Method GET -Headers $headers
    Write-Output "OK"
    Write-Output $r.id
    Write-Output $r.object
} catch {
    Write-Output "ERROR: $($_.Exception.Message)"
    Write-Output $_.ErrorDetails.Message
}
