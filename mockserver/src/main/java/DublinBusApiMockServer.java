import org.mockserver.configuration.ConfigurationProperties;
import org.mockserver.integration.ClientAndServer;

import static org.mockserver.integration.ClientAndServer.startClientAndServer;

public class DublinBusApiMockServer {

    public DublinBusApiMockServer() {
        ConfigurationProperties.enableCORSForAllResponses(true);
        ConfigurationProperties.corsAllowOrigin("*");
        ConfigurationProperties.persistExpectations(true);
        ConfigurationProperties.persistedExpectationsPath("MockServerInit.json");
        ConfigurationProperties.initializationJsonPath("init.json");
        ClientAndServer clientAndServer = startClientAndServer(1080);
    }

    public static void main(String[] args) {
        new DublinBusApiMockServer();
    }
}
