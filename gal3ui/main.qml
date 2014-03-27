import QtQuick 2.0
import "../gal3ui" as g

Rectangle {
	id: page
	width: 1024; height: 768
	color: "lightgray"

	Rectangle {
		id: toolbar
		width: parent.width; height: 30
		color: activePalette.window
		anchors.top: screen.top

		g.Button {
			anchors { left: parent.left; verticalCenter: parent.verticalCenter }
			text: "New Album"
			//onClicked: game.startNewGame(gameCanvas, dialog)
		}
	}

	g.Dialog {
		id: dialog
		anchors.centerIn: parent
		z: 100
	}

}
