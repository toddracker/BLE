package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/devicehive/IoT-framework/devicehive-cloud/conf"
	"github.com/devicehive/IoT-framework/devicehive-cloud/rest"
	"github.com/devicehive/IoT-framework/devicehive-cloud/ws"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
)

func wsImplementation(bus *dbus.Conn, config conf.Conf) {
	var conn *ws.Conn
	for {
		info, err := rest.GetDHServerInfo(config.URL)
		if err == nil {
			c := ws.New(info.WebSocketServerUrl, config.DeviceID, func(m map[string]interface{}) {
				//log.Printf("Successfully received command: %s", m)
				log.Printf("|| CLOUD received Command:%v\n", m)

				p := m["parameters"]
				params := ""

				if p != nil {
					b, err := json.Marshal(p)
					if err != nil {
						log.Panic(err)
					}

					params = string(b)
				}
				// log.Printf("Parameters: %v", params)
				bus.Emit("/com/devicehive/cloud",
					"com.devicehive.cloud.CommandReceived",
					uint32(m["id"].(float64)),
					m["command"].(string),
					params)
			})
			conn = &c

			if err == nil {
				break
			}
		}
		log.Print(err)
		time.Sleep(5 * time.Second)
	}

	w := NewDbusObjectWrapper(conn)
	go conn.Run(func() {
		conn.RegisterDevice(config.DeviceID, config.DeviceName)
		conn.Authenticate()
		conn.SubscribeCommands()
	})

	bus.Export(w, "/com/devicehive/cloud", DBusConnName)

	// Introspectable
	n := &introspect.Node{
		Name: "/com/devicehive/cloud",
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:    "com.devicehive.cloud",
				Methods: introspect.Methods(w),
			},
		},
	}

	bus.Export(introspect.NewIntrospectable(n), "/com/devicehive/cloud", "org.freedesktop.DBus.Introspectable")

	select {}

}
