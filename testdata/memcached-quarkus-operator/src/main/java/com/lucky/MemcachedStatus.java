package com.lucky;

import java.util.ArrayList;
import java.util.List;

public class MemcachedStatus {

    // Add Status information here

    private List<String> nodes;

    public List<String> getNodes() {
        if (nodes == null) {
            nodes = new ArrayList<>();
        }
        return nodes;
    }

    public void setNodes(List<String> nodes) {
        this.nodes = nodes;
    }
}
