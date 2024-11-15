{{- $eventDestinationRef := (index .SDKAPI.API.Shapes "EventDestinations").MemberRef }}
{{- $eventDestinationName := "EventDestination" }}

// set{{ $eventDestinationName }} sets a resource {{ $eventDestinationName }} type
// given the SDK type.
func setResource{{ $eventDestinationName }}(
    resp *svcsdk.{{ $eventDestinationName }},
) *svcapitypes.{{ $eventDestinationName }} {
    res := &svcapitypes.{{ $eventDestinationName }}{}

{{ GoCodeSetResourceForStruct .CRD "EventDestination" "res" $eventDestinationRef "resp" $eventDestinationRef 1 }}
    return res
}
