#include "ns3/core-module.h"
#include "ns3/network-module.h"
#include "ns3/wifi-module.h"
#include "ns3/internet-module.h"
#include "ns3/applications-module.h"
#include "ns3/mobility-module.h"
#include "ns3/flow-monitor-module.h"
#include "ns3/netanim-module.h" // Para NetAnim
#include <iomanip> // Para std::setprecision

#include <iostream>
#include <fstream>
#include <string>

using namespace ns3;

NS_LOG_COMPONENT_DEFINE("WifiInterference");

int main(int argc, char *argv[]) {
    bool verbose = false;
    bool tracing = true;
    double simulationTime = 30.0;
    uint32_t numClients = 3;

    CommandLine cmd;
    cmd.AddValue("verbose", "Enable verbose output", verbose);
    cmd.AddValue("tracing", "Enable pcap tracing", tracing);
    cmd.AddValue("simulationTime", "Simulation end time", simulationTime);
    cmd.AddValue("numClients", "Number of client nodes", numClients);
    cmd.Parse(argc, argv);

    if (verbose) {
        LogComponentEnable("UdpClient", LOG_LEVEL_INFO);
        LogComponentEnable("UdpServer", LOG_LEVEL_INFO);
        LogComponentEnable("WifiNetDevice", LOG_LEVEL_ALL);
        LogComponentEnable("MacLow", LOG_LEVEL_ALL);
        LogComponentEnable("Phy", LOG_LEVEL_ALL);
        LogComponentEnable("YansWifiPhy", LOG_LEVEL_ALL);
        LogComponentEnable("YansWifiChannel", LOG_LEVEL_ALL);
        LogComponentEnable("MobilityModel", LOG_LEVEL_INFO);
        LogComponentEnable("ArpCache", LOG_LEVEL_ALL);
    }

    // 1. Criar os Nós
    NodeContainer nodes;
    nodes.Create(numClients + 2); // 1 AP, numClients, 1 Interferer
    NodeContainer apNode = NodeContainer(nodes.Get(0));
    NodeContainer clientNodes;
    for (uint32_t i = 0; i < numClients; ++i) {
        clientNodes.Add(nodes.Get(i + 1));
    }
    NodeContainer interfererNode = NodeContainer(nodes.Get(numClients + 1));

    WifiHelper wifi;
    wifi.SetStandard(WIFI_STANDARD_80211b);
    YansWifiPhyHelper phy;
    YansWifiChannelHelper channel = YansWifiChannelHelper::Default();
    phy.SetChannel(channel.Create());
    wifi.SetRemoteStationManager("ns3::AarfWifiManager");
    // Configurar o MAC para o ponto de acesso (AP)
    Ssid ssid = Ssid("my-wifi-network");
    WifiMacHelper macAp;
    macAp.SetType("ns3::ApWifiMac",
                  "Ssid", SsidValue(ssid),
                  "BeaconInterval", TimeValue(Seconds(0.1024)));

    // Configurar o MAC para os clientes (STA - Station)
    WifiMacHelper macSta;
    macSta.SetType("ns3::StaWifiMac",
                   "Ssid", SsidValue(ssid),
                   "ActiveProbing", BooleanValue(false));

    // Criar os dispositivos Wi-Fi para cada nó
    NetDeviceContainer apDevice = wifi.Install(phy, macAp, apNode);
    NetDeviceContainer clientDevices = wifi.Install(phy, macSta, clientNodes);
    NetDeviceContainer interfererDevice = wifi.Install(phy, macSta, interfererNode);

    // 3. Configurar a Mobilidade
    MobilityHelper mobility;

    // Posição do Ponto de Acesso
    Ptr<ListPositionAllocator> positionAllocAp = CreateObject<ListPositionAllocator>();
    positionAllocAp->Add(Vector(0.0, 0.0, 0.0));
    mobility.SetPositionAllocator(positionAllocAp);
    mobility.Install(apNode);

    // Posição dos Clientes
    Ptr<ListPositionAllocator> positionAllocClients = CreateObject<ListPositionAllocator>();
    positionAllocClients->Add(Vector(5.0, 5.0, 0.0));
    positionAllocClients->Add(Vector(10.0, 5.0, 0.0));
    positionAllocClients->Add(Vector(5.0, 10.0, 0.0));
    mobility.SetPositionAllocator(positionAllocClients);
    mobility.Install(clientNodes);

    // Posição do Nó de Interferência
    Ptr<ListPositionAllocator> positionAllocInterferer = CreateObject<ListPositionAllocator>();
    positionAllocInterferer->Add(Vector(2.0, 2.0, 0.0));
    mobility.SetPositionAllocator(positionAllocInterferer);
    mobility.Install(interfererNode);

    // 4. Configurar a Pilha de Protocolos TCP/IP
    InternetStackHelper stack;
    stack.Install(nodes);

    Ipv4AddressHelper address;
    address.SetBase("10.1.1.0", "255.255.255.0");
    NetDeviceContainer allDevices;
    allDevices.Add(apDevice);
    allDevices.Add(clientDevices);
    allDevices.Add(interfererDevice);
    Ipv4InterfaceContainer interfaces = address.Assign(allDevices);
    Ipv4InterfaceContainer apIf;
    apIf.Add(interfaces.Get(0));
    Ipv4InterfaceContainer clientIfs;
    for (uint32_t i = 0; i < numClients; ++i) {
        clientIfs.Add(interfaces.Get(i + 1));
    }
    Ipv4InterfaceContainer interfererIf;
    interfererIf.Add(interfaces.Get(numClients + 1));

    // 5. Criar Aplicações

    // Servidor UDP no Ponto de Acesso
    uint16_t port = 9;
    UdpServerHelper server(port);
    ApplicationContainer serverApps = server.Install(apNode.Get(0));
    serverApps.Start(Seconds(1.0));
    serverApps.Stop(Seconds(simulationTime));

    // Cliente UDP nos Nós Clientes
    uint32_t packetSize = 1024;
    uint32_t numPackets = 1000;
    UdpClientHelper client(interfaces.GetAddress(0), port);
    client.SetAttribute("MaxPackets", UintegerValue(numPackets));
    client.SetAttribute("PacketSize", UintegerValue(packetSize));
    client.SetAttribute("Interval", TimeValue(Seconds(0.01))); // Enviar a cada 10ms

    ApplicationContainer clientApps;
    for (uint32_t i = 0; i < clientNodes.GetN(); ++i) {
        ApplicationContainer app = client.Install(clientNodes.Get(i));
        app.Start(Seconds(2.0 + i * 0.1)); // Iniciar os clientes em momentos ligeiramente diferentes
        app.Stop(Seconds(simulationTime));
        clientApps.Add(app);
    }

    // Gerador de Tráfego no Nó de Interferência
    uint16_t interfererPort = 10;
    UdpClientHelper interfererClient(interfaces.GetAddress(0), interfererPort); // Enviar para o servidor ou broadcast
    interfererClient.SetAttribute("MaxPackets", UintegerValue(0)); // Enviar continuamente
    interfererClient.SetAttribute("PacketSize", UintegerValue(512));
    interfererClient.SetAttribute("Interval", TimeValue(Seconds(0.005))); // Enviar a cada 5ms
    ApplicationContainer interfererApps = interfererClient.Install(interfererNode.Get(0));
    interfererApps.Start(Seconds(2.5));
    interfererApps.Stop(Seconds(simulationTime - 1.5));

    // 6. Monitoramento de Fluxo
    FlowMonitorHelper flowMonitor;
    Ptr<FlowMonitor> monitor = flowMonitor.InstallAll();

    // 7. Tracing (Opcional)
    if (tracing) {
        NetDeviceContainer pcapDevices;
        pcapDevices.Add(apDevice);
        pcapDevices.Add(clientDevices);
        pcapDevices.Add(interfererDevice);
        phy.EnablePcap("wifi-interference", pcapDevices);
    }

    // 8. NetAnim (Visualização com NetAnim)
    AnimationInterface anim("simulation_2.xml");
    anim.EnablePacketMetadata(true);
    
    // Nomeando os nós
    anim.UpdateNodeDescription(apNode.Get(0), "AP");
    anim.UpdateNodeColor(apNode.Get(0), 0, 255, 0); // Verde

    for (uint32_t i = 0; i < clientNodes.GetN(); ++i) {
        std::ostringstream clientName;
        clientName << "Client-" << i + 1;
        anim.UpdateNodeDescription(clientNodes.Get(i), clientName.str());
        anim.UpdateNodeColor(clientNodes.Get(i), 0, 0, 255); // Azul
    }

    anim.UpdateNodeDescription(interfererNode.Get(0), "Interferer");
    anim.UpdateNodeColor(interfererNode.Get(0), 255, 0, 0); // Vermelho

    // Posicionando os nós
    anim.SetConstantPosition(apNode.Get(0), 0.0, 0.0);

    std::vector<Vector> clientPositions = {
        Vector(5.0, 5.0, 0.0),
        Vector(10.0, 5.0, 0.0),
        Vector(5.0, 10.0, 0.0)
    };

    for (uint32_t i = 0; i < clientNodes.GetN(); ++i) {
        Vector pos = clientPositions[i % clientPositions.size()];
        anim.SetConstantPosition(clientNodes.Get(i), pos.x, pos.y);
    }

    anim.SetConstantPosition(interfererNode.Get(0), 2.0, 2.0);


    // 9. Executar a Simulação
    Simulator::Stop(Seconds(simulationTime));
    Simulator::Run();

    // 10. Analisar os Resultados
    monitor->CheckForLostPackets();
    FlowMonitor::FlowStatsContainer stats = monitor->GetFlowStats();

    for (auto i = stats.begin(); i != stats.end(); ++i) {
        std::cout << "Flow ID: " << i->first << std::endl;
        std::cout << "  Packets Sent: " << i->second.txPackets << std::endl;
        std::cout << "  Packets Received: " << i->second.rxPackets << std::endl;
        std::cout << "  Packets Lost: " << i->second.txPackets - i->second.rxPackets << std::endl;
        if (i->second.txPackets > 0) {
            double packetLossRate = (double)(i->second.txPackets - i->second.rxPackets) / i->second.txPackets * 100;
            std::cout << "  Packet Loss Rate: " << std::fixed << std::setprecision(2) << packetLossRate << "%" << std::endl;
        } else {
            std::cout << "  No packets sent." << std::endl;
        }
        std::cout << "-----------------------" << std::endl;
    }

    monitor->SerializeToXmlFile("simulation2.flowmon", false, false);

    Simulator::Destroy();
    return 0;
}