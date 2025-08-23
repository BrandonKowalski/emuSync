import {useEffect, useState} from 'react';
import './App.css';
import {models} from "../wailsjs/go/models";

import {ListDevices} from "../wailsjs/go/sync/Sync";
import Device = models.Device;

function App() {

    const [devices, setDevices] = useState<Device[]>([]);

    useEffect(() => {
        ListDevices().then((d) => {
            setDevices(d)
        })
    }, [devices]);


    return (
        <div id="App">
            <div className="devices">
                <h1>Devices:</h1>
                {devices.map((d) => {
                    return <div className="device" key={d.id}>
                        <span className="id">{d.id}</span><br />
                        {d.name} &nbsp;
                        <button className="btn" onClick={() => {
                        }}>Connect
                        </button>
                    </div>
                })}
            </div>
        </div>
    )
}

export default App
