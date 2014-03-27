import QtQuick 2.0

Rectangle {
	id: inputContainer

	function show(text) {
		container.opacity = 1;
	}

	function hide() {
		container.opacity = 0;
	}

	width: url.width + 20
	height: url.height + apiKey.height + urlLabel.height + apiKeyLabel.height + 20
	opacity: 0

	Text {
		id: urlLabel
		text: "Gallery URL"
	}

	TextInput {
		id: url
		anchors.centerIn: parent
		text: ""
	}   
	
	Text {
		id: apiKeyLabel
		text: "API Key"
	}

	TextInput {
		id: apiKey
		anchors.centerIn: parent
		text: ""
	}
	Behavior on opacity {
		NumberAnimation { properties:"opacity"; duration: 500 }
	}

	MouseArea {
		anchors.fill: parent
		onClicked: hide();
	}
}
