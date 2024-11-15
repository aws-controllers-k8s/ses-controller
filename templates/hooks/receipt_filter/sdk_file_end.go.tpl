{{- $receiptFilterRef := (index .SDKAPI.API.Shapes "ReceiptFilterList").MemberRef }}
{{- $receiptFilterName := "ReceiptFilter" }}

// set{{ $receiptFilterName }} sets a resource {{ $receiptFilterName }} type
// given the SDK type.
func setResource{{ $receiptFilterName }}(
    resp *svcsdk.{{ $receiptFilterName }},
) *svcapitypes.{{ $receiptFilterName }}_SDK {
    res := &svcapitypes.{{ $receiptFilterName }}_SDK{}

{{ GoCodeSetResourceForStruct .CRD "Rule" "res" $receiptFilterRef "resp" $receiptFilterRef 1 }}
    return res
}
