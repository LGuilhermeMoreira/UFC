#include "ns3/applications-module.h"
#include "ns3/core-module.h"
#include "ns3/internet-module.h"
#include "ns3/network-module.h"
#include "ns3/point-to-point-module.h"
#include "ns3/netanim-module.h"  

using namespace ns3;

NS_LOG_COMPONENT_DEFINE("MultiNodePacketDelivery");

uint32_t packetsSent = 0;
uint32_t packetsReceived = 0;

void TxCallback(Ptr<const Packet> packet)
{
    packetsSent++;
}

void RxCallback(Ptr<const Packet> packet)
{
    packetsReceived++;
}


int main(int argc, char* argv[])
{
    uint32_t numClients = 4;
    CommandLine cmd;
    cmd.AddValue("numClients", "Number of client nodes", numClients);
    cmd.Parse(argc, argv);

    Time::SetResolution(Time::NS);
    LogComponentEnable("UdpEchoClientApplication", LOG_LEVEL_INFO);
    LogComponentEnable("UdpEchoServerApplication", LOG_LEVEL_INFO);

    NodeContainer serverNode;
    serverNode.Create(1);

    NodeContainer clientNodes;
    clientNodes.Create(numClients);

    InternetStackHelper stack;
    stack.Install(serverNode);
    stack.Install(clientNodes);

    Ipv4AddressHelper address;
    PointToPointHelper p2p;
    p2p.SetDeviceAttribute("DataRate", StringValue("5Mbps"));
    p2p.SetChannelAttribute("Delay", StringValue("2ms"));

    UdpEchoServerHelper echoServer(9);
    ApplicationContainer serverApp = echoServer.Install(serverNode.Get(0));
    serverApp.Start(Seconds(1.0));
    serverApp.Stop(Seconds(20.0));

    for (uint32_t i = 0; i < numClients; ++i)
    {
        NetDeviceContainer devices = p2p.Install(clientNodes.Get(i), serverNode.Get(0));

        devices.Get(0)->TraceConnectWithoutContext("PhyTxEnd", MakeCallback(&TxCallback));
        devices.Get(1)->TraceConnectWithoutContext("PhyRxEnd", MakeCallback(&RxCallback));


        Ipv4InterfaceContainer interfaces;
        std::ostringstream subnet;
        subnet << "10.1." << i + 1 << ".0";
        address.SetBase(subnet.str().c_str(), "255.255.255.0");
        interfaces = address.Assign(devices);

        UdpEchoClientHelper echoClient(interfaces.GetAddress(1), 9);
        echoClient.SetAttribute("MaxPackets", UintegerValue(3));
        echoClient.SetAttribute("Interval", TimeValue(Seconds(1.0)));
        echoClient.SetAttribute("PacketSize", UintegerValue(1024));

        ApplicationContainer clientApp = echoClient.Install(clientNodes.Get(i));
        clientApp.Start(Seconds(2.0 + i)); 
        clientApp.Stop(Seconds(20.0));
    }

    AnimationInterface anim("simulation-1.xml");
    anim.SetConstantPosition(serverNode.Get(0), 50.0, 50.0);
    anim.UpdateNodeDescription(serverNode.Get(0), "Server");
    anim.UpdateNodeColor(serverNode.Get(0), 0, 255, 0); 

    for (uint32_t i = 0; i < numClients; ++i)
    {
        anim.SetConstantPosition(clientNodes.Get(i), 10.0 * (i + 1), 10.0);
        std::ostringstream label;
        label << "Client" << i;
        anim.UpdateNodeDescription(clientNodes.Get(i), label.str());
        anim.UpdateNodeColor(clientNodes.Get(i), 0, 0, 255); // azul
    }


    Simulator::Run();

    double deliveryRate = (packetsSent > 0) ? (double)packetsReceived / packetsSent * 100.0 : 0.0;

    std::cout << "Pacotes enviados: " << packetsSent << std::endl;
    std::cout << "Pacotes recebidos: " << packetsReceived << std::endl;
    std::cout << "Taxa de entrega: " << deliveryRate << "%" << std::endl;

    Simulator::Destroy();
    return 0;
}

