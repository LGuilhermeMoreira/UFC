#include <fstream>
#include <sstream> // Para construir chaves do mapa
#include <map>     // Para contar pacotes por fluxo
#include <ns3/csma-helper.h>
#include "ns3/mobility-module.h"
#include "ns3/nstime.h"
#include "ns3/core-module.h"
#include "ns3/network-module.h"
#include "ns3/internet-module.h"
#include "ns3/point-to-point-module.h"
#include "ns3/applications-module.h"
#include "ns3/ipv4-global-routing-helper.h"
#include "ns3/netanim-module.h"

#define TCP_SINK_PORT 9000
#define UDP_SINK_PORT 9001
#define DDOS_RATE "20480kb/s"
#define MAX_BULK_BYTES 150000
#define MAX_SIMULATION_TIME 50.0
#define NUMBER_OF_BOTS 20


using namespace ns3;
NS_LOG_COMPONENT_DEFINE("DDoSAttackAggregated");

std::ofstream csvFile;

// Estrutura para armazenar informações do fluxo
struct FlowData {
    uint16_t srcPort;
    uint16_t destPort;
    std::string transportLayer;
    std::string highestLayer;
    int target;
    int packetCount;
    uint64_t totalBytes;
};

// Mapa para contar pacotes por fluxo. A chave será uma string como "UDP-5000-9001"
std::map<std::string, FlowData> flowDataMap;
const double timeInterval = 1.0; // Intervalo de 1 segundo para agregação

void PacketCounter(std::string context, const Ipv4Header &ipHeader, Ptr<const Packet> packet, uint32_t interface) {
    // ... (lógica para extrair portas e protocolo permanece a mesma) ...
    uint8_t protocol = ipHeader.GetProtocol();
    uint16_t srcPort = 0;
    uint16_t dstPort = 0;
    std::string transportLayerStr;
    int isAttack = 0;

    if (protocol == 6) { /* TCP */ 
        TcpHeader tcpHeader;
        if (packet->PeekHeader(tcpHeader)) {
            srcPort = tcpHeader.GetSourcePort();
            dstPort = tcpHeader.GetDestinationPort();
            transportLayerStr = "TCP";
            isAttack = 0;
        }
    } else if (protocol == 17) { /* UDP */
        UdpHeader udpHeader;
        if (packet->PeekHeader(udpHeader)) {
            srcPort = udpHeader.GetSourcePort();
            dstPort = udpHeader.GetDestinationPort();
            transportLayerStr = "UDP";
            isAttack = 1;
        }
    } else { return; }

    std::stringstream flowKey;
    flowKey << transportLayerStr << "-" << srcPort << "-" << dstPort;

    if (flowDataMap.find(flowKey.str()) == flowDataMap.end()) {
        flowDataMap[flowKey.str()] = {srcPort, dstPort, transportLayerStr, transportLayerStr, isAttack, 0, 0};
    }
    
    flowDataMap[flowKey.str()].packetCount++;
    flowDataMap[flowKey.str()].totalBytes += packet->GetSize(); // <-- ATUALIZA O TOTAL DE BYTES
}

void WriteAggregatedData() {
    for (auto const& [key, data] : flowDataMap) {
        // Calcula o tamanho médio do pacote para este fluxo neste intervalo
        double avgPacketLength = 0;
        if (data.packetCount > 0) {
            avgPacketLength = static_cast<double>(data.totalBytes) / data.packetCount;
        }

        csvFile << data.srcPort << ","
                << data.destPort << ","
                << avgPacketLength << ","      // <-- USA O TAMANHO MÉDIO
                << data.packetCount << ","      // "Packets/Time"
                << data.highestLayer << ","
                << data.transportLayer << ","
                << data.target << std::endl;
    }
    flowDataMap.clear();
    if (Simulator::Now().GetSeconds() < MAX_SIMULATION_TIME - timeInterval) {
        Simulator::Schedule(Seconds(timeInterval), &WriteAggregatedData);
    }
}

int main(int argc, char *argv[]) {
    CommandLine cmd;
    cmd.Parse(argc, argv);

    Time::SetResolution(Time::NS);
    
    csvFile.open("attack_data_ml.csv");
    csvFile << "Source Port,Dest Port,Packet Length,Packets/Time,Highest Layer,Transport Layer,target" << std::endl;

    NodeContainer nodes;
    nodes.Create(3);
    NodeContainer botNodes;
    botNodes.Create(NUMBER_OF_BOTS);
    PointToPointHelper pp1, pp2;
    pp1.SetDeviceAttribute("DataRate", StringValue("100Mbps"));
    pp1.SetChannelAttribute("Delay", StringValue("1ms"));
    pp2.SetDeviceAttribute("DataRate", StringValue("100Mbps"));
    pp2.SetChannelAttribute("Delay", StringValue("1ms"));
    NetDeviceContainer d02, d12, botDeviceContainer[NUMBER_OF_BOTS];
    d02 = pp1.Install(nodes.Get(0), nodes.Get(1));
    d12 = pp1.Install(nodes.Get(1), nodes.Get(2));
    for (int i = 0; i < NUMBER_OF_BOTS; ++i)
    {
        botDeviceContainer[i] = pp2.Install(botNodes.Get(i), nodes.Get(1));
    }
    InternetStackHelper stack;
    stack.Install(nodes);
    stack.Install(botNodes);
    Ipv4AddressHelper ipv4_n;
    ipv4_n.SetBase("10.0.0.0", "255.255.255.252");
    Ipv4AddressHelper a02, a12;
    a02.SetBase("10.1.1.0", "255.255.255.0");
    a12.SetBase("10.1.2.0", "255.255.255.0");
    for (int j = 0; j < NUMBER_OF_BOTS; ++j)
    {
        ipv4_n.Assign(botDeviceContainer[j]);
        ipv4_n.NewNetwork();
    }
    Ipv4InterfaceContainer i02, i12;
    i02 = a02.Assign(d02);
    i12 = a12.Assign(d12);

    // ... (Configuração das aplicações é a mesma) ...
    OnOffHelper onoff("ns3::UdpSocketFactory", Address(InetSocketAddress(i12.GetAddress(1), UDP_SINK_PORT)));
    onoff.SetConstantRate(DataRate(DDOS_RATE));
    onoff.SetAttribute("OnTime", StringValue("ns3::ConstantRandomVariable[Constant=30]"));
    onoff.SetAttribute("OffTime", StringValue("ns3::ConstantRandomVariable[Constant=0]"));
    ApplicationContainer onOffApp[NUMBER_OF_BOTS];
    for (int k = 0; k < NUMBER_OF_BOTS; ++k)
    {
        onOffApp[k] = onoff.Install(botNodes.Get(k));
        onOffApp[k].Start(Seconds(0.0));
        onOffApp[k].Stop(Seconds(MAX_SIMULATION_TIME));
    }
    BulkSendHelper bulkSend("ns3::TcpSocketFactory", InetSocketAddress(i12.GetAddress(1), TCP_SINK_PORT));
    bulkSend.SetAttribute("MaxBytes", UintegerValue(MAX_BULK_BYTES));
    ApplicationContainer bulkSendApp = bulkSend.Install(nodes.Get(0));
    bulkSendApp.Start(Seconds(0.0));
    bulkSendApp.Stop(Seconds(MAX_SIMULATION_TIME - 10));
    PacketSinkHelper UDPsink("ns3::UdpSocketFactory", Address(InetSocketAddress(Ipv4Address::GetAny(), UDP_SINK_PORT)));
    ApplicationContainer UDPSinkApp = UDPsink.Install(nodes.Get(2));
    UDPSinkApp.Start(Seconds(0.0));
    UDPSinkApp.Stop(Seconds(MAX_SIMULATION_TIME));
    PacketSinkHelper TCPsink("ns3::TcpSocketFactory", InetSocketAddress(Ipv4Address::GetAny(), TCP_SINK_PORT));
    ApplicationContainer TCPSinkApp = TCPsink.Install(nodes.Get(2));
    TCPSinkApp.Start(Seconds(0.0));
    TCPSinkApp.Stop(Seconds(MAX_SIMULATION_TIME));

    Ipv4GlobalRoutingHelper::PopulateRoutingTables();

    // Conecta o trace source à função de CONTAGEM de pacotes
    Ptr<Ipv4> ipv4_server = nodes.Get(2)->GetObject<Ipv4>();
    ipv4_server->TraceConnect("LocalDeliver", "PacketCounterContext", MakeCallback(&PacketCounter));

    // Agenda a primeira chamada à função de ESCRITA agregada
    Simulator::Schedule(Seconds(timeInterval), &WriteAggregatedData);

    AnimationInterface anim("DDoSim.xml");
    MobilityHelper mobility;
    mobility.SetPositionAllocator("ns3::GridPositionAllocator", "MinX", DoubleValue(0.0), "MinY", DoubleValue(0.0), "DeltaX", DoubleValue(5.0), "DeltaY", DoubleValue(10.0), "GridWidth", UintegerValue(5), "LayoutType", StringValue("RowFirst"));
    mobility.SetMobilityModel("ns3::ConstantPositionMobilityModel");
    mobility.Install(nodes);
    mobility.Install(botNodes);
    ns3::AnimationInterface::SetConstantPosition(nodes.Get(0), 0, 0);
    ns3::AnimationInterface::SetConstantPosition(nodes.Get(1), 10, 10);
    ns3::AnimationInterface::SetConstantPosition(nodes.Get(2), 20, 10);
    uint32_t x_pos = 0;
    for (int l = 0; l < NUMBER_OF_BOTS; ++l)
    {
        ns3::AnimationInterface::SetConstantPosition(botNodes.Get(l), x_pos++, 30);
    }
    
    Simulator::Run();

    // Fechar o arquivo CSV
    csvFile.close();

    Simulator::Destroy();
    return 0;
}