	input.SetTemplate(&svcsdk.Template{
		TemplateName: desired.ko.Spec.Name,
		HtmlPart: desired.ko.Spec.HTMLPart,
		TextPart: desired.ko.Spec.TextPart,
		SubjectPart: desired.ko.Spec.SubjectPart,
	})
