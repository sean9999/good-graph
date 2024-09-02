const nodeGraph = document.getElementById('graph');

//  buttons
const addNode1 = document.getElementById('btnAddNode1');
const addNode2 = document.getElementById('btnAddNode2');
const lonely = document.getElementById('btnLonely');
const emitPartices = document.getElementById('btnEmitPartices');
const debug = document.getElementById('btnDebug');
const egalitarian = document.getElementById('btnEgalitarian');
const hier = document.getElementById('btnHier');
const marco = document.getElementById('btnMarco');
const killYourself = document.getElementById('btnKillYourself');
const reconnect = document.getElementById('btnReconnect');
const colourNode = document.getElementById('btnColourNode');

const btn = {
    addNode1,
    addNode2,
    lonely,
    emitPartices,
    debug,
    egalitarian,
    hier,
    marco,
    killYourself,
    reconnect,
    colourNode
};

export { btn, nodeGraph };

