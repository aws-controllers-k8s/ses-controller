{{- $receiptRuleRef := (index .SDKAPI.API.Shapes "ReceiptRulesList").MemberRef }}
{{- $receiptRuleName := "ReceiptRule" }}

// set{{ $receiptRuleName }} sets a resource {{ $receiptRuleName }} type
// given the SDK type.
func setResource{{ $receiptRuleName }}(
    resp *svcsdk.{{ $receiptRuleName }},
) *svcapitypes.{{ $receiptRuleName }}_SDK {
    res := &svcapitypes.{{ $receiptRuleName }}_SDK{}

{{ GoCodeSetResourceForStruct .CRD "Rule" "res" $receiptRuleRef "resp" $receiptRuleRef 1 }}
    return res
}
